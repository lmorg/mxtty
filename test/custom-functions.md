# Custom Functions

This is an example of a markdown table with custom formula functions

| Character | Show | On-screen age | Town/City |
| --- | --- | --- | --- |
| Bart Simpson | The Simpsons | 10 | Springfield |
| Lisa Simpson | The Simpsons | 8 | Springfield |
| Stewie Griffin | Family Guy | 1 | Quahog, Rhode Island |
| Peter Griffin | Family Guy | 44 | Quahog, Rhode Island |
| Eric Cartman | South Park | 10 | South Park, Colorado |
| Timmy Turner | The Fairly OddParents | 10 | Dimmsdale, California |
| Jimmy Neutron | Jimmy Neutron: Boy Genius | 11 | Retroville, Texas |
| Dipper Pines | Gravity Falls | 12 | Gravity Falls, Oregon |
| Arnold Shortman | Hey Arnold! | 9 | Hillwood |
| Ben Tennyson | Ben 10 | 10 | Bellwood |
| **Longest:** | **Longest:** | **Largest:** | **Longest:** |
| =longest(A2:A11) | =longest(B2:B11) | =largest(C2:C11) | =longest(D2:D11)

# Functions

## longest

Python function to return the longest cell value

```py
import sys

args = sys.argv[1:]
longest = max(args, key=len)
print(longest)
```

## largest

Bash function to return the largest cell value

```bash
highest=$1

for arg in "$@"; do
	((arg > highest)) && highest=$arg
done

echo $highest
```