// this is hell

use std::{
    process,
    error::Error,
    path::Path,
};

use super::functions::json::parse_json_file;

const RESOURCES : &str = "resources";

pub const OLD_EXECUTABLE : &str = "mozuli.exe.old";
pub const CURRENT_EXECUTABLE : &str = "mozuli.exe";

pub const CURRENT_VERSION : &str = "1.0.0";
pub const API_VERSION : &str = "v1";

pub static DEVELOPMENT_MODE : bool = Path::new("../development").exists();

pub static PROCESS_ID : String = (process::id()).to_string();

pub static RESOURCES_DATA : String = RESOURCES.to_owned() + "/data";

// FIXME
pub static APPLICATION_DATA : Result<serde_json::Value, Box<dyn Error>> = parse_json_file((RESOURCES_DATA.to_owned() + "/application.json").as_str());
pub static USER_DATA : Result<serde_json::Value, Box<dyn Error>> = parse_json_file((RESOURCES_DATA.to_owned() + "/user.json").as_str());