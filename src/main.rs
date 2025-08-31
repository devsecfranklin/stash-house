use std::fs::File;
use std::io::{self, prelude::*, BufReader};
//use std::io::{self, BufRead};
use std::path::Path;

fn main() -> io::Result<()> {
    if let Ok(lines) = read_lines("./files.txt") {
        for line in lines {
            let mut rs:bool=true;
            if let Ok(loc) = line {
                rs = Path::new(&loc).exists();
    
                if rs == true{
                    println!("{:?} exists", loc);
                }
                else{
                    println!("Does not exist: {:?}", loc);
                }  
            }
        }
    }

    let file = File::open("dirs.txt")?;
    let reader = BufReader::new(file);
    for line in reader.lines() {
        //println!("{}", line?);
        let path: String = line.unwrap().as_str().unwrap().to_string();
        //for entry in std::fs::read_dir(line).expect("Unable to list") {
        //    let entry = entry.expect("unable to get entry");
        //    println!("{}", entry.path().display());
        //}
    }
    Ok(())
}

fn read_lines<P>(filename: P) -> io::Result<io::Lines<io::BufReader<File>>>
where P: AsRef<Path>, {
    let file = File::open(filename)?;
    Ok(io::BufReader::new(file).lines())
}
