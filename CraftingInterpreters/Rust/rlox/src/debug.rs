use crate::{chunk::{Chunk, OpCodes}, value::LoxValue};

pub fn disassemble_chunk(chunk: Chunk, name: String) {
    println!("{}", name);
    let mut offset = 0;
    while offset < chunk.code.len() {
        let code = chunk.code[offset];
        let opcode = OpCodes::from_int(code);
        print!("{:04} ", offset);
        match opcode {
            OpCodes::OpReturn => offset += simple_instruction(String::from("OP_RETURN")),
            OpCodes::OpConstant => offset += constant_instruction(String::from("OP_CONSTANT"), &chunk, offset)
        }
    }
}

fn simple_instruction(name: String) -> usize {
    println!("{}", name);
    return 1;
}

fn constant_instruction(name: String, chunk: &Chunk, offset: usize) -> usize {
    let constant : LoxValue = chunk.constants[usize::from(chunk.code[offset + 1])];
    println!("{}, {:.2}", name, constant);
    return 2;
}