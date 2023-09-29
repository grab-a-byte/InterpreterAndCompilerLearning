const std = @import("std");
const Chunk = @import("chunk.zig").Chunk;
const OpCode = @import("chunk.zig").OpCode;

pub const InterpretResult = enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR,
};

pub const VM = struct {
    const Self = @This();

    chunk: *Chunk,
    ip: u8,

    pub fn init(chunk: *Chunk) VM {
        return VM{ .chunk = chunk, .ip = 0 };
    }

    pub fn interpret(self: *Self, isDebug: bool) InterpretResult {
        return self.run(isDebug);
    }

    fn run(self: *Self, isDebug: bool) InterpretResult {
        while (true) {
            const instruction: OpCode = @enumFromInt(self.chunk.code.items[self.ip]);
            self.ip += 1;
            switch (instruction) {
                OpCode.OP_RETURN => {
                    return InterpretResult.INTERPRET_OK;
                },
                OpCode.OP_CONSTANT => {
                    const constant = self.chunk.constants.items[self.chunk.code.items[self.ip]];
                    self.ip += 1;
                    if (isDebug) {
                        std.debug.prin("{}\n", .{constant});
                    }
                },
            }
        }
    }

    pub fn free(self: *Self) void {
        self.chunk.free();
    }
};
