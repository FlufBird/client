use std::{
    process,
    error::Error,
    path::Path,
};

use super::functions::json::parse_json_file;

pub struct GlobalVariables {
    pub current_version : &'static str,
    pub api_version : &'static str,

    pub old_executable : &'static str,
    pub current_executable : &'static str,

    pub development_mode : bool,

    pub process_id : String,

    pub resources : &'static str,

    pub resources_data : String,
    pub resources_languages : String,

    pub application_data : Result<serde_json::Value, Box<dyn Error>>,
    pub user_data : Result<serde_json::Value, Box<dyn Error>>,
}

pub fn set() -> GlobalVariables {
    const RESOURCES : &str = "resources";

    let resources_data : String = RESOURCES.to_owned() + "/data";
    let resources_languages : String = RESOURCES.to_owned() + "/languages";

    let global_variables = GlobalVariables {
        current_version : "1.0.0",
        api_version : "v1",

        old_executable : "mozuli.exe.old",
        current_executable : "mozuli.exe",

        development_mode : Path::new("../development").exists(),

        process_id : (process::id()).to_string(),

        resources : RESOURCES,

        resources_data : (&resources_data).to_owned(),
        resources_languages : (&resources_languages).to_owned(),

        application_data : parse_json_file(((&resources_data).to_owned() + "/application.json").as_str()),
        user_data : parse_json_file(((&resources_languages).to_owned() + "/user.json").as_str()),
    };

    global_variables
}