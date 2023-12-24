const std = @import("std");
const Value = @import("value.zig").Value;
const allocator = std.heap.page_allocator;

pub const OpCode = enum(u8) {
    OP_RETURN,
    OP_CONSTANT,
    OP_NEGATE,
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
};

pub const Chunk = struct {
    const Self = @This();
    code: std.ArrayList(u8),
    constants: std.ArrayList(Value),
    lines: std.ArrayList(u32),

    pub fn init() Chunk {
        return Chunk{ .lines = std.ArrayList(u32).init(allocator), .code = std.ArrayList(u8).init(allocator), .constants = std.ArrayList(Value).init(allocator) };
    }

    pub fn addConstant(self: *Self, value: Value) !usize {
        try self.constants.append(value);
        return self.constants.items.len - 1;
    }

    pub fn writeOp(self: *Self, op: OpCode, line: u32) !void {
        try self.write(@intFromEnum(op), line);
    }

    pub fn write(self: *Self, value: u8, line: u32) !void {
        try self.code.append(value);
        try self.lines.append(line);
    }

    pub fn free(self: *Chunk) void {
        self.constants.deinit();
        self.code.deinit();
        self.lines.deinit();
    }
};
