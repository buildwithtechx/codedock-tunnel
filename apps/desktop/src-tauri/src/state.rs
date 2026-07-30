use std::process::Child;
use std::sync::Mutex;

pub struct TunnelState {
    pub(crate) child: Mutex<Option<Child>>,
}

impl Default for TunnelState {
    fn default() -> Self {
        Self {
            child: Mutex::new(None),
        }
    }
}
