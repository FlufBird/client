use std::{
    error::Error,
    fs::File,
    io::{
        Read,
        Write,

        BufReader,
    },
};

pub fn open_file(path : &str) -> Result<File, Box<dyn Error>> {
    Ok(File::open(path) ?)
}

pub fn read_file(path : &str) -> Result<String, Box<dyn Error>> {
    let _file = open_file(path) ?;
    let mut reader = BufReader::new(_file);
    let mut content = String::new();

    reader.read_to_string(&mut content) ?;

    Ok(content)
}

pub fn write_file(path : &str, content : String) {
    let __file = open_file(path);
    let mut _file;

    match __file {
        Ok(___file) => _file = ___file,
        Err(_) => return,
    }

    match _file.write_all(content.as_bytes()) {
        Ok(_) => (),
        Err(_) => (),
    } // handling errors so rust shuts up
}