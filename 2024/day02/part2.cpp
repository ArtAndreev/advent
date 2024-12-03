#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>

int main() {
    std::ifstream input("input.txt");
    std::string line;
    size_t res = 0;
    while (std::getline(input, line)) {
        std::istringstream iss(line);

        std::vector<int> report;
        int num;
        while (iss >> num) {
            report.push_back(num);
        }

        if (report.size() == 0) {
            ++res;
            continue;
        }

        // Increasing.
        int wrong_count = 0;
        int prev = report.front();
        for (size_t i = 1; i < report.size(); ++i) {
            int num = report[i];
            if (num <= prev || num - prev > 3) {
                ++wrong_count;
                if (wrong_count > 1) {
                    break;
                }
            } else {
                prev = num;
            }
        }

        if (wrong_count < 2) {
            ++res;
            continue;
        }

        // Decreasing.
        wrong_count = 0;
        prev = report.front();
        for (size_t i = 1; i < report.size(); ++i) {
            int num = report[i];
            if (num >= prev || prev - num > 3) {
                ++wrong_count;
                if (wrong_count > 1) {
                    break;
                }
            } else {
                prev = num;
            }
        }

        if (wrong_count < 2) {
            ++res;
            continue;
        }

        // Increasing from back.
        wrong_count = 0;
        prev = report.back();
        for (ssize_t i = report.size() - 2; i >= 0; --i) {
            int num = report[i];
            if (num <= prev || num - prev > 3) {
                ++wrong_count;
                if (wrong_count > 1) {
                    break;
                }
            } else {
                prev = num;
            }
        }

        if (wrong_count < 2) {
            ++res;
            continue;
        }

        // Decreasing from back.
        wrong_count = 0;
        prev = report.back();
        for (ssize_t i = report.size() - 2; i >= 0; --i) {
            int num = report[i];
            if (num >= prev || prev - num > 3) {
                ++wrong_count;
                if (wrong_count > 1) {
                    break;
                }
            } else {
                prev = num;
            }
        }

        if (wrong_count < 2) {
            ++res;
        }
    }

    std::cout << res << std::endl;

    return 0;
}
