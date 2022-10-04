extern crate serde_json;
extern crate reqwest;

use std::{
    process,
    error,
    string,
    thread,
    path,
    fs,
};

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

fn delete_old_version() {
    let _path = path::Path::new("old_mozuli.exe");

    if _path.exists() {
        fs::remove_file(_path).ok();
    }
}

fn check_updates() {
}

fn update() {
}

fn display_dialog(title : &str, message : &str) {
}

fn display_critical_error(message : &str) {
    display_dialog("Critical Error", message);

    process::exit(1);
}

fn main() {
    const API_VERSION : &str = "v1";

    let development_mode = path::Path::new("../development").exists();

    let application_data = parse_json_file("resources/data/application.json");
    let user_data = parse_json_file("resources/data/user.json");

    match application_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't parse application data."),
    };

    match user_data {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't parse user data."),
    };

    let server = ((match development_mode {
        true => "http://localhost:5000",
        false => "https://mozuli.deta.dev",
    }).to_owned() + "/api/") + API_VERSION;

    let http = reqwest::Client::builder().build();

    match http {
        Ok(_) => (),
        Err(_) => display_critical_error("Couldn't build HTTP client."),
    }

    delete_old_version();

    // TODO spawn thread to check for updates every 30 seconds
}