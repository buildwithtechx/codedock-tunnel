#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod commands;
mod state;

use commands::{tunnel_start, tunnel_status, tunnel_stop, tunnel_version};
use state::TunnelState;

fn main() {
    tauri::Builder::default()
        .manage(TunnelState::default())
        .invoke_handler(tauri::generate_handler![
            tunnel_start,
            tunnel_stop,
            tunnel_status,
            tunnel_version
        ])
        .run(tauri::generate_context!())
        .expect("error while running Codedock Tunnel desktop application");
}
