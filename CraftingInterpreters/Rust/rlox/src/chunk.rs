use crate::value::LoxValue;

pub enum OpCodes {
    OpReturn = 0,
    OpConstant = 1,
}

impl OpCodes {
    pub fn from_int(value: u8) -> OpCodes {
        match value {
            x if x == Self::OpReturn as u8 => Self::OpReturn,
            x if x == Self::OpConstant as u8 => Self::OpConstant,
            _ => panic!("unkown opcode"),
        }
    }
}

pub struct Chunk {
    pub code: Vec<u8>,
    pub constants: Vec<LoxValue>,
}

impl Chunk {
    pub fn init() -> Chunk {
        return Chunk {
            code: Vec::new(),
            constants: Vec::new(),
        };
    }

    pub fn write_op(&mut self, code: OpCodes) {
        self.write(code as u8);
    }

    pub fn add_constant(&mut self, value: LoxValue) -> usize {
        self.constants.push(value);
        return self.constants.len() - 1;
    }

    pub fn write(&mut self, item: u8) {
        self.code.push(item);
    }
}
