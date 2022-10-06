extern crate serde_json;
extern crate tauri;

use std::{
    process,
    error,
    string,
    thread,
    time::Duration,
    path,
    fs,
};

use tauri::{
    api::{
        http,
    },
};

fn delete_old_version(old_executable : &str) {
    let _path = path::Path::new(old_executable);

    if _path.exists() {
        fs::remove_file(_path).ok();
    }
}

fn check_instances(application_data : &serde_json::Value, process_id : &String) {
    if application_data["processId"] == process_id.to_string() {
        process::exit(0);
    }
}

fn write_instance(application_data : &serde_json::Value, process_id : &String) {
    // TODO: edit "processId" in application_data
}

fn read_file(path : &str) -> Result<string::String, Box<dyn error::Error>> {
    Ok(fs::read_to_string(path) ?)
}

fn parse_json(data : &str) -> Result<serde_json::Value, Box<dyn error::Error>> {
    Ok(serde_json::from_str(data) ?)
}

fn parse_json_file(path : &str) -> Result<serde_json::Value, Box<dyn error::Error>> {
    let content = read_file(path) ?;
    let data = parse_json(content.as_str()) ?;

    Ok(data)
}

async fn send_request(requests : &http::Client, method : &str, url : &str) -> Result<Result<http::Response, tauri::api::Error>, Box<dyn error::Error>> {
    let request = http::HttpRequestBuilder::new(method, url) ?;

    let response = requests.send(request).await;

    Ok(response)
}

async fn check_updates(application_data : &serde_json::Value, requests : &http::Client, api_update : &String) -> Result<bool, Box<dyn error::Error>> {
    let latest_version;

    match send_request(requests, "get", &(api_update.to_owned() + "/latest_version")).await {
        Ok(response) => latest_version = response,
        Err(_) => (), // TODO: return error
    }

    Ok(false)
}

fn update() {
}

fn display_dialog(title : &str, message : &str) {
    // title form: [title] - Mozuli
}

fn display_critical_error(message : &str) {
    display_dialog("Critical Error", message);

    process::exit(1);
}

fn main() {
    const API_VERSION : &str = "v1";

    const OLD_EXECUTABLE : &str = "mozuli.exe.old";
    const CURRENT_EXECUTABLE : &str = "mozuli.exe";

    let process_id = (process::id()).to_string();

    let development_mode = path::Path::new("../development").exists();

    let _application_data = parse_json_file("resources/data/application.json");
    let _user_data = parse_json_file("resources/data/user.json");

    let application_data;
    let user_data;

    match _application_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve application data."),
    };

    application_data = _application_data.unwrap();

    match _user_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't retrieve user data."),
    };

    user_data = _user_data.unwrap();

    let server = (match development_mode {
        true => "http://localhost:5000",
        false => "https://mozuli.deta.dev",
    }).to_owned();

    let api = (&server).to_owned() + "/api/" + API_VERSION;

    let api_update = api + "/update";

    let _requests = http::ClientBuilder::new()
        .connect_timeout(Duration::from_secs(10))
        .max_redirections(0)

        .build();
    let requests;

    match _requests {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't build HTTP client."),
    }

    check_instances(&application_data, &process_id);
    write_instance(&application_data, &process_id);

    requests = _requests.unwrap();

    delete_old_version(OLD_EXECUTABLE);

    let mut check_updates_thread_running = true;

    // thread::spawn(move || async { // FIXME
    //     let interval = Duration::from_secs(30);

    //     loop {
    //         if !check_updates_thread_running {
    //             break;
    //         }

    //         match check_updates(&application_data, &requests, &api_update).await {
    //             Ok(result) => {
    //                 match result {
    //                     false => (),
    //                     true => break,
    //                 }
    //             },
    //             Err(_) => (),
    //         }

    //         thread::sleep(interval);
    //     }

    //     update();
    // });

    check_updates_thread_running = false;
}