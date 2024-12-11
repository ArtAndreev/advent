#include <bitset>
#include <fstream>
#include <iostream>
#include <sstream>
#include <string>
#include <unordered_map>
#include <vector>

class BigInt {
  public:
    BigInt operator+(size_t other) {
        bool has_overflow = false;
        auto num_bucket_it = num.rbegin();
        while (other != 0) {
            uint16_t res_bucket = other & 0b11111111;
            if (has_overflow) {
                ++res_bucket;
            }

            if (num_bucket_it != num.rend()) {
                uint16_t num_bucket = *num_bucket_it;
                res_bucket += num_bucket;
                *num_bucket_it = res_bucket & 0b11111111;
                ++num_bucket_it;
            } else {
                uint8_t new_bucket = res_bucket & 0b11111111;
                num.insert(num.begin(), new_bucket);
                num_bucket_it = num.rend();
            }

            has_overflow = res_bucket & 0b100000000;
            other >>= 8;
        }
        while (has_overflow) {
            if (num_bucket_it == num.rend()) {
                num.insert(num.begin(), 1);
                break;
            }

            uint16_t num_bucket = *num_bucket_it;
            ++num_bucket;
            *num_bucket_it = num_bucket & 0b11111111;
            ++num_bucket_it;

            has_overflow = num_bucket & 0b100000000;
        }

        return *this;
    }

    std::string to_binary_string() const {
        if (num.empty()) {
            return "0";
        }

        std::stringstream iss;
        for (const auto& bucket : num) {
            iss << std::bitset<8>(bucket);
        }

        return iss.str();
    }

  private:
    // friend BigInt operator+(const BigInt& other) {}
    // friend BigInt operator+=(const BigInt& other) {}

    std::vector<uint8_t> num;
};

int main() {
    std::ifstream input("input.txt");

    std::vector<size_t> stones;
    size_t num;
    while (input >> num) {
        stones.push_back(num);
    }

    std::unordered_map<size_t, size_t> stones_with_count;
    for (size_t stone : stones) {
        ++stones_with_count[stone];
    }

    for (size_t i = 0; i < 75; ++i) {
        std::unordered_map<size_t, size_t> new_stones_with_count;
        for (auto it = stones_with_count.begin(); it != stones_with_count.end(); ++it) {
            size_t stone = it->first;
            size_t count = it->second;
            if (stone == 0) {
                new_stones_with_count[1] += count;
            } else if (std::string stone_str = std::to_string(stone); stone_str.size() % 2 == 0) {
                size_t first = std::stoll(stone_str.substr(0, stone_str.size() / 2));
                size_t second = std::stoll(stone_str.substr(stone_str.size() / 2));

                new_stones_with_count[first] += count;
                new_stones_with_count[second] += count;
            } else {
                // no overflow there, tested.
                new_stones_with_count[stone * 2024] += count;
            }
        }

        stones_with_count = std::move(new_stones_with_count);
    }

    // Should have used double...
    BigInt res;
    for (auto it = stones_with_count.begin(); it != stones_with_count.end(); ++it) {
        res = res + it->second;
    }
    // TODO: decimal int.
    std::cout << res.to_binary_string() << std::endl;

    return 0;
}
