use std::error::Error;

use super::_file::{
    read_file,
    write_file,
};

pub fn string_to_json(data : &str) -> Result<serde_json::Value, Box<dyn Error>> {
    Ok(serde_json::from_str(data) ?)
}

pub fn json_to_string(data : serde_json::Value) -> Result<String, Box<dyn Error>> {
    Ok(serde_json::to_string(&data) ?)
}

pub fn parse_json_file(path : &str) -> Result<serde_json::Value, Box<dyn Error>> {
    Ok(string_to_json(read_file(path) ?.as_str()) ?)
}

pub fn update_json_file(path : &str, data : &serde_json::Value) {
    let _content = json_to_string(data.to_owned());
    let content;

    match _content {
        Ok(__content) => content = __content,
        Err(_) => return,
    }

    write_file(path, content);
}