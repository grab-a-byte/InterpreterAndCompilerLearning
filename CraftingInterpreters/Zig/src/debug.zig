const std = @import("std");
const OpCodes = @import("chunk.zig").OpCode;
const Chunk = @import("chunk.zig").Chunk;
const Value = @import("value.zig").Value;
const print = std.debug.print;
const allocator = std.heap.page_allocator;

pub fn disassembleChunk(input: *Chunk, name: []const u8) void {
    std.debug.print("{any}", .{input.constants.items});
    std.debug.print("{any}", .{input.code.items});
    print("{s}\n", .{name});
    const items = input.code.items;
    var offset: usize = 0;

    while (offset < items.len) {
        const byte = items[offset];
        print("{d:0>4} ", .{offset});
        const toPrint: []const u8 =
            if (offset > 0 and input.lines.items[offset] == input.lines.items[offset - 1]) "   | " else std.fmt.allocPrint(allocator, "{d:0>4} ", .{input.lines.items[offset]}) catch "Unable to get line number";

        print("{s}", .{toPrint});
        offset += dissassembleInstruction(input, offset, byte);
    }
}

pub fn dissassembleInstruction(input: *Chunk, offset: usize, byte: u8) u8 {
    const enumVal: OpCodes = @enumFromInt(byte);
    switch (enumVal) {
        OpCodes.OP_CONSTANT => return constantInstruction("OP_CONSTANT", input, offset),
        OpCodes.OP_RETURN => return simpleInstruction("OP_RETURN"),
        OpCodes.OP_NEGATE => return simpleInstruction("OP_NEGATE"),
        OpCodes.OP_ADD => return simpleInstruction("OP_ADD"),
        OpCodes.OP_SUBTRACT => return simpleInstruction("OP_SUBTRACT"),
        OpCodes.OP_MULTIPLY => return simpleInstruction("OP_MULTIPLY"),
        OpCodes.OP_DIVIDE => return simpleInstruction("OP_DIVIDE"),
    }
}

fn simpleInstruction(comptime name: []const u8) u8 {
    print("{s} \n", .{name});
    return 1;
}

fn constantInstruction(comptime name: []const u8, chunk: *Chunk, offset: usize) u8 {
    const constant = chunk.constants.items[chunk.code.items[offset + 1]];
    print("{s} : {d:.2} \n", .{ name, constant });
    return 2;
}

pub fn printValue(value: Value) void {
    std.debug.print("{d}", .{value});
}
