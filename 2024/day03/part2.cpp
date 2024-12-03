#include <cassert>
#include <cstdlib>
#include <fstream>
#include <iostream>

enum class State {
    initial,
    reading_do_dont_cmd,
    reading_do_cmd,
    reading_dont_cmd,
    reading_mul_cmd,
    reading_first_arg,
    reading_second_arg,
};

static const char* mul_cmd = "mul(";
static const size_t mul_cmd_len = strlen(mul_cmd);

char get_expected_mul_cmd_char(int read_mul_cmd_progress) {
    assert(read_mul_cmd_progress < mul_cmd_len);
    return mul_cmd[read_mul_cmd_progress];
}

static const char* do_cmd = "do()";
static const size_t do_cmd_len = strlen(do_cmd);

char get_expected_do_cmd_char(int read_do_cmd_progress) {
    assert(read_do_cmd_progress < do_cmd_len);
    return do_cmd[read_do_cmd_progress];
}

static const char* dont_cmd = "don't()";
static const size_t dont_cmd_len = strlen(dont_cmd);

char get_expected_dont_cmd_char(int read_dont_cmd_progress) {
    assert(read_dont_cmd_progress < dont_cmd_len);
    return dont_cmd[read_dont_cmd_progress];
}

static const int common_do_len = 2;

int char_to_int(char ch) {
    return ch - '0';
}

int main() {
    std::ifstream input("input.txt");

    size_t res;
    bool is_mul_enabled = true;
    int read_cmd_progress = 0;
    int first_arg = 0;
    int second_arg = 0;
    State state = State::initial;
    char ch;
    while (input.get(ch)) {
        switch (state) {
        case State::initial:
            if (ch == get_expected_do_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                state = State::reading_do_dont_cmd;
            } else if (is_mul_enabled && ch == get_expected_mul_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                state = State::reading_mul_cmd;
            } else {
                read_cmd_progress = 0;
            }
            break;

        case State::reading_do_dont_cmd:
            if (ch == get_expected_do_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                if (read_cmd_progress > common_do_len) {
                    state = State::reading_do_cmd;
                }
            } else if (ch == get_expected_dont_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                if (read_cmd_progress > common_do_len) {
                    state = State::reading_dont_cmd;
                }
            } else {
                read_cmd_progress = 0;
                state = State::initial;
            }
            break;

        case State::reading_do_cmd:
            if (ch == get_expected_do_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                if (read_cmd_progress == do_cmd_len) {
                    is_mul_enabled = true;
                    read_cmd_progress = 0;
                    state = State::initial;
                }
            } else {
                read_cmd_progress = 0;
                state = State::initial;
            }
            break;

        case State::reading_dont_cmd:
            if (ch == get_expected_dont_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                if (read_cmd_progress == dont_cmd_len) {
                    is_mul_enabled = false;
                    read_cmd_progress = 0;
                    state = State::initial;
                }
            } else {
                read_cmd_progress = 0;
                state = State::initial;
            }
            break;

        case State::reading_mul_cmd:
            if (ch == get_expected_mul_cmd_char(read_cmd_progress)) {
                ++read_cmd_progress;
                if (read_cmd_progress == mul_cmd_len) {
                    read_cmd_progress = 0;
                    state = State::reading_first_arg;
                }
            } else {
                read_cmd_progress = 0;
                state = State::initial;
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
                state = State::initial;
            }
            break;

        case State::reading_second_arg:
            // first char cannot be zero while another not zero,
            // also arg is guaranteed to be length > 0 and < 4, just add to arg.
            if (ch == ')') {
                res += first_arg * second_arg;
                first_arg = 0;
                second_arg = 0;
                state = State::initial;
            } else if (ch >= '0' && ch <= '9') {
                second_arg = second_arg * 10 + char_to_int(ch);
            } else {
                first_arg = 0;
                second_arg = 0;
                state = State::initial;
            }
            break;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
