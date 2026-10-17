#!/bin/bash
# Converts Morphik's JSON OpenAPI 3.1 spec to 3.0.3 for kin-openapi.
# Nullable anyOf branches use nullable:true; a sole non-null branch is inlined,
# with references wrapped in allOf so nullable and sibling metadata survive.
# Numeric exclusive bounds become minimum/maximum plus 3.0 boolean flags.
# All operations and non-null union alternatives remain in the document.
set -euo pipefail
trap 'on_error "$?" "$LINENO"' ERR

LOG_FILE="${LOG_FILE:-/tmp/$(basename "$0" .sh).log}"
exec > >(tee -a "$LOG_FILE") 2>&1

SPEC_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)/api"
readonly SPEC_DIR

on_error() {
	local exit_code="$1"
	local line="$2"

	printf '{"level":"ERROR","file":"%s","line":%d,"msg":"command failed","exit_code":%d}\n' \
		"${BASH_SOURCE[1]}" "$line" "$exit_code" >&2
}

log() {
	local level="$1"
	shift

	jq -nc \
		--arg time "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" \
		--arg level "$level" \
		--arg file "${BASH_SOURCE[1]##*/}" \
		--arg line "${BASH_LINENO[0]}" \
		--arg func "${FUNCNAME[1]:-main}" \
		--arg msg "$*" \
		'{time:$time,level:$level,file:$file,line:($line|tonumber),func:$func,msg:$msg}' >&2
}

# shellcheck disable=SC2016 # jq variables and $ref are not shell expansions.
readonly PATCH_FILTER='
walk(
  if type == "object" then
    (if (.anyOf? | type) == "array" and any(.anyOf[]; .type? == "null") then
      .anyOf |= map(select(.type? != "null"))
      | .nullable = true
      | if (.anyOf | length) == 1 then
          .anyOf[0] as $branch | del(.anyOf)
          | if ($branch | has("$ref")) then .allOf = [$branch]
            else $branch + . end
        else . end
    else . end)
    | if (.exclusiveMinimum? | type) == "number" then
        .minimum = .exclusiveMinimum | .exclusiveMinimum = true
      else . end
    | if (.exclusiveMaximum? | type) == "number" then
        .maximum = .exclusiveMaximum | .exclusiveMaximum = true
      else . end
  else . end
) | .openapi = "3.0.3"
'

main() {
	local tmp
	tmp="$(mktemp "${SPEC_DIR}/.openapi.fixed.XXXXXX")"
	# shellcheck disable=SC2064 # expand tmp before the local leaves scope.
	trap "rm -f -- '${tmp}'" EXIT

	jq "$PATCH_FILTER" "${SPEC_DIR}/openapi.yaml" >"$tmp"
	chmod 0644 "$tmp"
	mv -- "$tmp" "${SPEC_DIR}/openapi.fixed.yaml"
	log INFO "converted OpenAPI 3.1 nullable unions and exclusive bounds"
}

main "$@"
