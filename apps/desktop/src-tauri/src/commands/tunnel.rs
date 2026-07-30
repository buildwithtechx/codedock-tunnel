use serde::Serialize;
use std::process::Command;
use tauri::State;

use crate::state::TunnelState;

#[derive(Serialize)]
pub struct TunnelProcess {
    pub pid: u32,
    pub status: String,
    pub exit_code: Option<i32>,
}

#[tauri::command]
pub fn tunnel_start(
    state: State<'_, TunnelState>,
    port: u16,
    protocol: String,
    subdomain: Option<String>,
) -> Result<TunnelProcess, String> {
    validate_protocol(&protocol)?;
    let mut child_slot = state
        .child
        .lock()
        .map_err(|_| "tunnel state is unavailable")?;
    if let Some(child) = child_slot.as_mut() {
        if child
            .try_wait()
            .map_err(|error| error.to_string())?
            .is_none()
        {
            return Err("a tunnel process is already running".to_string());
        }
    }
    *child_slot = None;
    let mut command = Command::new(cli_path());
    command
        .arg("open")
        .arg("--port")
        .arg(port.to_string())
        .arg("--protocol")
        .arg(protocol);
    if let Some(value) = subdomain.filter(|value| !value.trim().is_empty()) {
        command.arg("--subdomain").arg(value);
    }
    let child = command
        .spawn()
        .map_err(|error| format!("start tunnel CLI: {error}"))?;
    let pid = child.id();
    *child_slot = Some(child);
    Ok(TunnelProcess {
        pid,
        status: "running".to_string(),
        exit_code: None,
    })
}

#[tauri::command]
pub fn tunnel_stop(state: State<'_, TunnelState>) -> Result<TunnelProcess, String> {
    let mut child_slot = state
        .child
        .lock()
        .map_err(|_| "tunnel state is unavailable")?;
    let Some(mut child) = child_slot.take() else {
        return Ok(stopped_process(None));
    };
    child
        .kill()
        .map_err(|error| format!("stop tunnel CLI: {error}"))?;
    let status = child
        .wait()
        .map_err(|error| format!("wait for tunnel CLI: {error}"))?;
    Ok(stopped_process(status.code()))
}

#[tauri::command]
pub fn tunnel_status(state: State<'_, TunnelState>) -> Result<TunnelProcess, String> {
    let mut child_slot = state
        .child
        .lock()
        .map_err(|_| "tunnel state is unavailable")?;
    let Some(child) = child_slot.as_mut() else {
        return Ok(stopped_process(None));
    };
    match child.try_wait().map_err(|error| error.to_string())? {
        Some(status) => {
            let result = stopped_process(status.code());
            *child_slot = None;
            Ok(result)
        }
        None => Ok(TunnelProcess {
            pid: child.id(),
            status: "running".to_string(),
            exit_code: None,
        }),
    }
}

#[tauri::command]
pub fn tunnel_version() -> Result<String, String> {
    let output = Command::new(cli_path())
        .arg("version")
        .output()
        .map_err(|error| format!("run tunnel CLI: {error}"))?;
    if !output.status.success() {
        return Err(String::from_utf8_lossy(&output.stderr).trim().to_string());
    }
    Ok(String::from_utf8_lossy(&output.stdout).trim().to_string())
}

fn cli_path() -> String {
    std::env::var("CODEDOCK_TUNNEL_CLI_PATH").unwrap_or_else(|_| "codedock-tunnel".to_string())
}

fn validate_protocol(protocol: &str) -> Result<(), String> {
    match protocol {
        "http" | "https" | "tcp" | "udp" => Ok(()),
        _ => Err(format!("unsupported tunnel protocol: {protocol}")),
    }
}

fn stopped_process(exit_code: Option<i32>) -> TunnelProcess {
    TunnelProcess {
        pid: 0,
        status: "stopped".to_string(),
        exit_code,
    }
}
