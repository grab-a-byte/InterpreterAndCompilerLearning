mod chunk;
mod debug;
mod value;

fn main() {
    let mut c = chunk::Chunk::init();
    c.write_op(chunk::OpCodes::OpReturn);
    let index = c.add_constant(1.2);
    c.write_op(chunk::OpCodes::OpConstant);
    c.write(index.try_into().unwrap());
    debug::disassemble_chunk(c, String::from("test chunk"))
}
