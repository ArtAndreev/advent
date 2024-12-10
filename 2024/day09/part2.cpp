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
    int offset_from_end = 0;
    while (true) {
        auto end_file_it = disk.end() - 1 - offset_from_end;
        if (end_file_it->type == NodeType::space) {
            ++offset_from_end;
            continue;
        }

        auto begin_space_it = disk.begin() + offset_from_begin;
        if (begin_space_it->type == NodeType::file) {
            ++offset_from_begin;
            continue;
        }

        if (begin_space_it > end_file_it) {
            break;
        }

        while (begin_space_it < end_file_it && (begin_space_it->type == NodeType::file ||
                                                end_file_it->count > begin_space_it->count)) {
            ++begin_space_it;
        }
        if (begin_space_it == end_file_it) {
            ++offset_from_end;
            continue;
        }

        int count = end_file_it->count;
        Node new_file_node{end_file_it->file_id, count, NodeType::file};
        begin_space_it = disk.insert(begin_space_it, new_file_node) + 1;

        if (begin_space_it->count == count) {
            disk.erase(begin_space_it);
        } else {
            begin_space_it->count -= count;
        }

        end_file_it = disk.end() - 1 - offset_from_end;
        end_file_it = disk.erase(end_file_it);
        Node new_space_node{-1, count, NodeType::space};
        disk.insert(end_file_it, new_space_node);
        ++offset_from_end;
    }

    size_t res = 0;
    size_t i = 0;
    for (const auto& node : disk) {
        if (node.type == NodeType::space) {
            i += node.count;
            continue;
        }
        for (size_t j = 0; j < node.count; ++j) {
            res += i * node.file_id;
            ++i;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
