extern crate serde_json;

use std::{
    process,
    error,
    string,
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
    let data = parse_json(content.as_str()) ?; // FIXME

    Ok(data)
}

fn delete_old_version() {
    let _path = path::Path::new("old_mozuli.exe");

    if _path.exists() {
        fs::remove_file(_path).ok();
    }
}

fn check_updates() { // TODO
}

fn update() {} // TODO

fn display_dialog(title : &str, message : &str) { // TODO

}

fn display_critical_error(error : &str) { // TODO
    process::exit(1);
}

fn main() {
    const API_VERSION : &str = "v1";

    let development_mode = path::Path::new("../development").exists();

    let application_data = match parse_json_file("resources/data/application.json") {
        Ok(data) => data,
        Err(error) => {
            display_critical_error("Couldn't parse application data!");

            error // FIXME
        },
    };

    let server = ((match development_mode {
        true => "http://localhost:5000",
        false => "https://mozuli.deta.dev",
    }).to_owned() + "/api/") + API_VERSION;

    delete_old_version();
    check_updates();
}