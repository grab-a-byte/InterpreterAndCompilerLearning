const std = @import("std");
const Value = @import("value.zig").Value;
const allocator = std.heap.page_allocator;

pub const OpCode = enum(u8) { OP_RETURN, OP_CONSTANT };

pub const Chunk = struct {
    code: std.ArrayList(u8),
    constants: std.ArrayList(Value),
    lines: std.ArrayList(u32),

    pub fn init() Chunk {
        return Chunk{ .lines = std.ArrayList(u32).init(allocator), .code = std.ArrayList(u8).init(allocator), .constants = std.ArrayList(Value).init(allocator) };
    }

    pub fn addConstant(self: *Chunk, value: Value) !usize {
        try self.constants.append(value);
        return self.constants.items.len - 1;
    }

    pub fn writeOp(self: *Chunk, op: OpCode, line: u32) !void {
        try self.write(@intFromEnum(op), line);
    }

    pub fn write(self: *Chunk, value: u8, line: u32) !void {
        try self.code.append(value);
        try self.lines.append(line);
    }

    pub fn free(self: *Chunk) void {
        self.constants.deinit();
        self.code.deinit();
        self.lines.deinit();
    }
};
