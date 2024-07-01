# 统计代码行数

git log --author="Decision" --since="2024-05-28" --until="2024-07-01" --pretty=tformat: --numstat -- . ":(exclude)static" ":(exclude)example" ":(exclude)go.mod" | awk '{ add += $1; subs += $2; loc += $1 - $2 } END { printf "added lines: %s, removed lines: %s, total lines: %s\n", add, subs, loc }'
