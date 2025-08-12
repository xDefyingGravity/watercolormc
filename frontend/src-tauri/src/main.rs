#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use std::process::{Command, Child};
use std::sync::{Arc, Mutex};
use tauri::{Manager, WindowEvent};

fn main() {
    tauri::Builder::default()
        .setup(|app| {
            let resource_dir = app.path().resource_dir().unwrap();

            #[cfg(target_os = "windows")]
            let backend_path = resource_dir.join("resources/app.exe");
            #[cfg(not(target_os = "windows"))]
            let backend_path = resource_dir.join("resources/app");

            let child = Command::new(backend_path)
                .spawn()
                .expect("failed to launch backend");

            app.manage(Arc::new(Mutex::new(child)));

            Ok(())
        })
        .on_window_event(|window, event| {
            if let WindowEvent::CloseRequested { .. } = event {
                let app_handle = window.app_handle();

                let child_arc = window.state::<Arc<Mutex<Child>>>().inner().clone();

                if let Ok(mut child) = child_arc.lock() {
                    let _ = child.kill();
                    let _ = child.wait();
                }

                app_handle.exit(0);
            }
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri app");
}
