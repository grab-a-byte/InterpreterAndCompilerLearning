const std = @import("std");
const Chunk = @import("chunk.zig").Chunk;
const OpCodes = @import("chunk.zig").OpCode;
const debug = @import("debug.zig");
const VM = @import("vm.zig").VM;
const InterpretResult = @import("vm.zig").InterpretResult;

const print = std.debug.print;

pub fn main() !void {
    var newChunk = Chunk.init();
    const constIndex = try newChunk.addConstant(1.2);
    try newChunk.writeOp(OpCodes.OP_CONSTANT, 123);
    try newChunk.write(@as(u8, @truncate(constIndex)), 345);
    try newChunk.writeOp(OpCodes.OP_NEGATE, 123);
    try newChunk.writeOp(OpCodes.OP_RETURN, 123);
    _ = debug.disassembleChunk(&newChunk, "test chunk");

    std.debug.print("\n", .{});

    var vm = VM.init(&newChunk);
    _ = vm.interpret(true);

    newChunk.free();
}
