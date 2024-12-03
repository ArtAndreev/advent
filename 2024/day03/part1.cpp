#include <cassert>
#include <cstdlib>
#include <fstream>
#include <iostream>

enum class State {
    reading_cmd,
    reading_first_arg,
    reading_second_arg,
};

static const char* mul_cmd = "mul(";
static const size_t mul_cmd_len = strlen(mul_cmd);

char get_expected_cmd_char(int read_mul_cmd_progress) {
    assert(read_mul_cmd_progress < mul_cmd_len);
    return mul_cmd[read_mul_cmd_progress];
}

int char_to_int(char ch) {
    return ch - '0';
}

int main() {
    std::ifstream input("input.txt");

    size_t res;
    int read_mul_cmd_progress = 0;
    int first_arg = 0;
    int second_arg = 0;
    State state = State::reading_cmd;
    char ch;
    while (input.get(ch)) {
        switch (state) {
        case State::reading_cmd:
            if (ch == get_expected_cmd_char(read_mul_cmd_progress)) {
                ++read_mul_cmd_progress;
                if (read_mul_cmd_progress == mul_cmd_len) {
                    read_mul_cmd_progress = 0;
                    state = State::reading_first_arg;
                }
            } else {
                read_mul_cmd_progress = 0;
            }
            break;

        case State::reading_first_arg:
            // first char cannot be zero while another not zero,
            // also arg is guaranteed to be length > 0 and < 4, just add to arg.
            if (ch == ',') {
                state = State::reading_second_arg;
            } else if (ch >= '0' && ch <= '9') {
                first_arg = first_arg * 10 + char_to_int(ch);
            } else {
                first_arg = 0;
                state = State::reading_cmd;
            }
            break;

        case State::reading_second_arg:
            // first char cannot be zero while another not zero,
            // also arg is guaranteed to be length > 0 and < 4, just add to arg.
            if (ch == ')') {
                res += first_arg * second_arg;
                first_arg = 0;
                second_arg = 0;
                state = State::reading_cmd;
            } else if (ch >= '0' && ch <= '9') {
                second_arg = second_arg * 10 + char_to_int(ch);
            } else {
                first_arg = 0;
                second_arg = 0;
                state = State::reading_cmd;
            }
            break;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
