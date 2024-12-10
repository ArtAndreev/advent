#include <algorithm>
#include <deque>
#include <fstream>
#include <iostream>

enum class NodeType {
    file,
    space,
};

struct Node {
    int file_id;
    int count;
    NodeType type;
};

int main() {
    std::ifstream input("input.txt");

    std::deque<Node> disk;
    int file_id = 0;
    NodeType type = NodeType::file;
    char ch;
    while (input.get(ch) && ch != '\n') {
        int count = ch - '0';

        if (count != 0) {
            Node new_node{type == NodeType::file ? file_id : -1, count, type};
            disk.push_back(new_node);
        }

        if (type == NodeType::file) {
            ++file_id;
        }
        type = type == NodeType::file ? NodeType::space : NodeType::file;
    }

    int offset_from_begin = 0;
    while (true) {
        auto begin_space_it = disk.begin() + offset_from_begin;
        if (begin_space_it == disk.end()) {
            break;
        }
        if (begin_space_it->type == NodeType::file) {
            ++offset_from_begin;
            continue;
        }

        auto end_file_it = disk.end() - 1;
        if (end_file_it->type == NodeType::space) {
            disk.erase(end_file_it);
            continue;
        }

        int count = std::min(begin_space_it->count, end_file_it->count);
        Node new_file_node{end_file_it->file_id, count, NodeType::file};
        begin_space_it = disk.insert(begin_space_it, new_file_node) + 1;
        ++offset_from_begin;

        if (begin_space_it->count == count) {
            disk.erase(begin_space_it);
        } else {
            begin_space_it->count -= count;
        }

        end_file_it = disk.end() - 1;
        if (end_file_it->count == count) {
            disk.erase(end_file_it);
        } else {
            end_file_it->count -= count;
        }
    }

    size_t res = 0;
    size_t i = 0;
    for (const auto& node : disk) {
        for (size_t j = 0; j < node.count; ++j) {
            res += i * node.file_id;
            ++i;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
