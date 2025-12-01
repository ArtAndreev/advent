use std::fs::File;
use std::io::{BufRead, BufReader};

pub fn part_one() {
    let mut dial = 50;
    let file = File::open("input.txt").unwrap();
    let reader = BufReader::new(file);
    let mut zero_count = 0;
    for line in reader.lines() {
        let line = line.unwrap();
        let (direction, count) = line.split_at(1);
        let count = count.parse::<i32>().unwrap();

        dial = match direction {
            "L" => dial - count,
            "R" => dial + count,
            _ => panic!("unknown direction"),
        } % 100;
        if dial < 0 {
            dial = 100 - -dial;
        }

        if dial == 0 {
            zero_count += 1;
        }
    }

    println!("{}", zero_count);
}
