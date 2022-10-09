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

use super::global::variables::{
    OLD_EXECUTABLE,
    CURRENT_EXECUTABLE,

    CURRENT_VERSION,
    API_VERSION,

    DEVELOPMENT_MODE,

    PROCESS_ID,

    RESOURCES_DATA,

    APPLICATION_DATA,
    USER_DATA,
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

    match send_request(requests, "get", &(api_update.to_owned() + "/latest_version")) {
        Ok(response) => latest_version = response,
        Err(error) => {
            return Err(error);
        },
    }

    // TODO: check if latest_version["latestVersion"] == "currentVersion", if they arent, update()

    Ok(false) // placeholder
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

fn display_dialog(title : &str, message : &str) {
    // title form: {title} - Mozuli
}

fn display_critical_error(message : &str) {
    display_dialog("Critical Error", message);

    exit(1);
}

pub fn backend() {
    match APPLICATION_DATA {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve application data."),
    };

    let application_data = APPLICATION_DATA.unwrap();

    match USER_DATA {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve user data."),
    };

    let user_data = USER_DATA.unwrap();

    let server = (match DEVELOPMENT_MODE {
        true => "http://localhost:5000",
        false => "https://mozuli.deta.dev",
    }).to_owned();

    let api = server.to_owned() + "/api/" + API_VERSION;

    let api_update = api + "/update";

    let _requests = reqwest::blocking::ClientBuilder::new()
        .timeout(Duration::from_secs(10))

        .build();
    let requests;

    match _requests {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't build HTTP client."),
    }

    check_instances(application_data["processId"], &PROCESS_ID);
    write_instance(&application_data, &PROCESS_ID);

    requests = _requests.unwrap();

    delete_old_version(OLD_EXECUTABLE);

    thread::spawn(move || {
        let interval = Duration::from_secs(30);

        loop {
            match check_updates(CURRENT_VERSION.to_string(), &requests, &api_update) {
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
    });
}