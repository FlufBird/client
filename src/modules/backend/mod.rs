use std::{
    process::{
        exit,
    },
    error::Error,
    time::Duration,
    thread,
    path::Path,
    fs::remove_file,
};

use super::{
    global::variables,
    frontend::frontend,
};

use tauri::api::dialog::{
    MessageDialogBuilder,
    MessageDialogKind,
    MessageDialogButtons,
};

fn delete_old_version(old_executable : &str) {
    let _path = Path::new(old_executable); // to differentiate between the variable and the library

    if _path.exists() {
        match remove_file(_path) {
            Ok(_) => (),
            Err(_) => (),
        } // handling errors so rust shuts up
    }
}

fn check_instances(latest_process_id : String, current_process_id : &String) {
    if latest_process_id == current_process_id.to_string() {
        exit(0);
    }
}

fn write_instance(application_data : &serde_json::Value, process_id : &String) {
    // TODO: modify application_data["processId"] to process_id and save data
}

fn check_updates(current_version : String, requests : &reqwest::blocking::Client, api_update : &String) -> Result<bool, Box<dyn Error>> {
    let latest_version;

    match send_request(requests, "get", (api_update.to_owned() + "/latest_version").as_str()) {
        Ok(response) => latest_version = response,
        Err(error) => return Err(error),
    }

    Ok(false) // TODO: this is a placeholder, check if latest_version["latestVersion"] == "currentVersion"
}

fn update_checker(current_version : String, requests : &reqwest::blocking::Client, url : &String) {
    let interval = Duration::from_secs(30);

    loop {
        match check_updates(current_version, requests, url) {
            Ok(result) => {
                match result {
                    false => (),
                    true => {
                        thread::spawn(update);

                        break;
                    },
                }
            },
            Err(_) => (),
        }

        thread::sleep(interval);
    }
}

fn update() {
}

fn send_request(requests : &reqwest::blocking::Client, method : &str, url : &str) -> Result<reqwest::blocking::Response, Box<dyn Error>> {
    let function = match method {
        "get" => requests.get(url),
        "post" => requests.post(url),
        "put" => requests.put(url),
        "delete" => requests.delete(url),
        _ => requests.get(url), // this case shouldnt happen, but it needed to be handled and be the same type as others anyway
    };

    Ok(function.send() ?)
}

fn display_critical_error(message : &str) { // this function is only used for displaying critical errors, dialogs for the frontend are shown inside the webview
    MessageDialogBuilder::new("Critical Error - Mozuli", message)
        .kind(MessageDialogKind::Error)
        .buttons(MessageDialogButtons::Ok)
        .show(|_| {
            exit(1);
        }
    );

    exit(1);
}

pub fn backend() {
    let global_variables = variables::set();

    match global_variables.application_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve application data."),
    }

    match global_variables.user_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve user data."),
    }

    // FIXME: unwrap without making a new variable since it has to be global

    let server = match global_variables.development_mode {
        true => "http://localhost:5000",
        false => "https://mozuli.deta.dev",
    };

    let api = server.to_owned() + "/api" + "/v" + global_variables.api_version;

    let api_update = api + "/update";

    let _requests = reqwest::blocking::ClientBuilder::new()
        .timeout(Duration::from_secs(10))

        .build();
    let requests;

    match _requests {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't build HTTP client."),
    }

    check_instances(global_variables.application_data["processId"].to_string(), &global_variables.process_id);
    write_instance(&global_variables.application_data, &global_variables.process_id);

    requests = _requests.unwrap();

    delete_old_version(global_variables.old_executable);

    frontend(global_variables);

    thread::spawn(|| update_checker(global_variables.current_version.to_string(), &requests, &api_update));
}