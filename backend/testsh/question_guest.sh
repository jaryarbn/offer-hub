#!/usr/bin/env bash
# 题库/题目游客接口一键回归脚本。
#
# 依赖：bash、curl、python3；后端服务和 MongoDB 测试数据需已就绪。
# 用法：
#   ./testsh/question_guest.sh
#   BASE_URL=http://127.0.0.1:8180 ./testsh/question_guest.sh
#
# 本脚本仅覆盖当前 router.go 注册的游客 GET 接口，不携带 Token。
# 后续增加鉴权后，可复制本文件为 question_auth.sh，统一在 curl_args
# 中加入 Authorization 请求头，再补充收藏、点赞等登录态专属接口；不要
# 把登录态断言混入本脚本，以便游客回归始终可以独立执行。

set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8180}"
BASE_URL="${BASE_URL%/}"
CURL_TIMEOUT="${CURL_TIMEOUT:-10}"

if [[ -t 1 && -z "${NO_COLOR:-}" ]]; then
  GREEN='\033[0;32m'
  RED='\033[0;31m'
  YELLOW='\033[0;33m'
  RESET='\033[0m'
else
  GREEN=''
  RED=''
  YELLOW=''
  RESET=''
fi

passed=0
failed=0
response_status=""
response_error=""
validation_output=""
work_dir="$(mktemp -d "${TMPDIR:-/tmp}/offer-hub-question-test.XXXXXX")"
response_body="${work_dir}/response.json"
trap 'rm -rf "$work_dir"' EXIT

pass() {
  passed=$((passed + 1))
  printf '%b[PASS]%b %s\n' "$GREEN" "$RESET" "$1"
}

fail() {
  failed=$((failed + 1))
  printf '%b[FAIL]%b %s' "$RED" "$RESET" "$1"
  if [[ -n "${2:-}" ]]; then
    printf ' - %s' "$2"
  fi
  printf '\n'
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '%b[FAIL]%b 缺少依赖：%s\n' "$RED" "$RESET" "$1" >&2
    exit 2
  fi
}

request_get() {
  local path="$1"
  shift
  response_error=""
  : >"$response_body"

  local curl_args=(
    --silent
    --show-error
    --location
    --connect-timeout "$CURL_TIMEOUT"
    --max-time "$CURL_TIMEOUT"
    --output "$response_body"
    --write-out '%{http_code}'
  )

  if ! response_status="$(curl "${curl_args[@]}" --get "${BASE_URL}${path}" "$@" 2>"${work_dir}/curl.err")"; then
    response_error="$(<"${work_dir}/curl.err")"
    response_status="000"
  fi
}

# validate_json 会同时检查 HTTP 200、统一响应 code=0/msg=success 和接口结构。
# 成功时，all/list 输出 bank_id，list 输出 question_id，供后续依赖步骤使用。
validate_json() {
  local mode="$1"
  shift
  python3 - "$response_status" "$response_body" "$mode" "$@" <<'PY'
import json
import sys

http_status, body_path, mode, *args = sys.argv[1:]

if http_status != "200":
    raise SystemExit(f"HTTP 状态应为 200，实际为 {http_status}")

try:
    with open(body_path, "r", encoding="utf-8") as file:
        payload = json.load(file)
except (OSError, json.JSONDecodeError) as error:
    raise SystemExit(f"响应不是有效 JSON：{error}")

def require(condition, message):
    if not condition:
        raise SystemExit(message)

require(isinstance(payload, dict), "响应根节点应为对象")

if mode == "missing_detail":
    require(payload.get("code") == 400, f"业务 code 应为 400，实际为 {payload.get('code')!r}")
    require(payload.get("msg") == "question_id is required", f"msg 不符合约定：{payload.get('msg')!r}")
    require(payload.get("data") is None, "参数错误时 data 应为 null")
    raise SystemExit(0)

require(payload.get("code") == 0, f"业务 code 应为 0，实际为 {payload.get('code')!r}")
require(payload.get("msg") == "success", f"msg 应为 success，实际为 {payload.get('msg')!r}")

data = payload.get("data")

if mode == "all_list":
    require(isinstance(data, list) and data, "data 应为非空职位方向列表")
    for job in data:
        require(isinstance(job.get("job_name"), str) and job["job_name"], "缺少 job_name")
        require(isinstance(job.get("series_list"), list), "series_list 应为数组")
        for series in job["series_list"]:
            require(isinstance(series.get("series_id"), str) and series["series_id"], "缺少 series_id")
            require(isinstance(series.get("bank_list"), list), "bank_list 应为数组")
            for bank in series["bank_list"]:
                require(isinstance(bank.get("count"), int) and bank["count"] >= 0, "题库 count 应为非负整数")
                if bank["count"] > 0:
                    print(bank["bank_id"])
                    raise SystemExit(0)
    raise SystemExit("未找到包含正常题目的题库")

if mode == "question_list":
    bank_id = args[0]
    require(isinstance(data, dict), "data 应为对象")
    require(isinstance(data.get("total"), int) and data["total"] > 0, "total 应大于 0")
    questions = data.get("list")
    require(isinstance(questions, list) and questions, "题目 list 应为非空数组")
    first = questions[0]
    require(bank_id in first.get("bank_list", []), "列表结果未命中请求的 bank_id")
    require(first.get("status") == 1, "游客列表只能返回正常题目")
    require(first.get("user_tag") == 0, "游客 user_tag 应为 0")
    require(first.get("user_liked") is False, "游客 user_liked 应为 false")
    question_id = first.get("question_id")
    require(isinstance(question_id, str) and question_id, "缺少 question_id")
    print(question_id)

elif mode == "meta_list":
    require(isinstance(data, dict), "data 应为对象")
    require(isinstance(data.get("total"), int) and data["total"] > 0, "total 应大于 0")
    items = data.get("list")
    require(isinstance(items, list) and items, "元数据 list 应为非空数组")
    require(set(items[0]) == {"question_id", "title"}, "元数据项只能包含 question_id、title")

elif mode == "detail":
    expected_id = args[0]
    require(isinstance(data, dict), "详情 data 应为对象")
    require(data.get("question_id") == expected_id, "详情 question_id 与列表抽取值不一致")
    require(data.get("status") == 1, "详情应为正常题目")
    require(isinstance(data.get("content"), str) and data["content"], "详情 content 不能为空")
    require(len(data["content"]) <= 150, "游客详情 content 最多返回前 150 个字符")
    require(data.get("analysis_content") == "", "游客 analysis_content 应为空字符串")
    require(data.get("user_tag") == 0, "游客 user_tag 应为 0")
    require(data.get("user_liked") is False, "游客 user_liked 应为 false")

elif mode == "hot_list":
    limit = int(args[0])
    require(isinstance(data, dict), "data 应为对象")
    items = data.get("list")
    require(isinstance(items, list) and 0 < len(items) <= limit, f"热门列表数量应在 1 到 {limit} 之间")
    expected_fields = {"question_id", "bank_list", "title", "view_count"}
    require(all(set(item) == expected_fields for item in items), "热门题目字段不符合接口约定")

else:
    raise SystemExit(f"未知校验模式：{mode}")
PY
}

run_validation() {
  local step="$1"
  local mode="$2"
  shift 2

  local output
  if output="$(validate_json "$mode" "$@" 2>&1)"; then
    pass "$step"
    validation_output="$output"
    return 0
  fi

  local detail="$output"
  if [[ -n "$response_error" ]]; then
    detail="$response_error"
  elif [[ -z "$detail" ]]; then
    detail="响应校验失败"
  fi
  fail "$step" "$detail"
  validation_output=""
  return 0
}

require_command curl
require_command python3

printf '%b题目模块游客回归：%s%b\n\n' "$YELLOW" "$BASE_URL" "$RESET"

request_get "/api/v1/question/all/list" \
  --data-urlencode "job_name=后端开发"
run_validation "GET /api/v1/question/all/list" all_list
bank_id="$validation_output"

question_id=""
if [[ -n "$bank_id" ]]; then
  request_get "/api/v1/question/list" \
    --data-urlencode "bank_id=$bank_id" \
    --data-urlencode "job_name=后端开发" \
    --data-urlencode "page=1" \
    --data-urlencode "page_size=5" \
    --data-urlencode "sort_by=order" \
    --data-urlencode "sort_order=asc"
  run_validation "GET /api/v1/question/list" question_list "$bank_id"
  question_id="$validation_output"
else
  fail "GET /api/v1/question/list" "上一步没有取得 bank_id"
fi

if [[ -n "$bank_id" ]]; then
  request_get "/api/v1/question/meta/list" \
    --data-urlencode "bank_id=$bank_id" \
    --data-urlencode "page=1" \
    --data-urlencode "page_size=5"
  run_validation "GET /api/v1/question/meta/list" meta_list
else
  fail "GET /api/v1/question/meta/list" "题库列表步骤没有取得 bank_id"
fi

if [[ -n "$question_id" ]]; then
  request_get "/api/v1/question/detail" \
    --data-urlencode "question_id=$question_id"
  run_validation "GET /api/v1/question/detail" detail "$question_id"
else
  fail "GET /api/v1/question/detail" "题目列表步骤没有取得 question_id"
fi

request_get "/api/v1/question/hot/list" \
  --data-urlencode "job_name=后端开发" \
  --data-urlencode "limit=3"
run_validation "GET /api/v1/question/hot/list" hot_list 3

request_get "/api/v1/question/detail"
run_validation "GET /api/v1/question/detail 缺少 question_id" missing_detail

printf '\n------------------------------\n'
printf '通过：%b%d%b  失败：%b%d%b  总计：%d\n' \
  "$GREEN" "$passed" "$RESET" "$RED" "$failed" "$RESET" "$((passed + failed))"

if ((failed > 0)); then
  exit 1
fi

exit 0
