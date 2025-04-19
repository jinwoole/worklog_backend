#!/usr/bin/env bash
set -e

API_URL="http://localhost:8080/api"
EMAIL="user@example.com"
PASSWORD="password123"

echo "1) 회원가입"
curl -s -X POST "$API_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}" | jq .
echo -e "\n"

echo "2) 로그인 및 토큰 획득"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\"}")
echo "$LOGIN_RESPONSE" | jq .
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r .token)
echo "TOKEN=$TOKEN"
echo -e "\n"

echo "3) 내 정보 조회 (Get Me)"
curl -s "$API_URL/me" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

echo "4) 오늘 업무 로그 작성 (Create WorkLog)"
curl -s -X POST "$API_URL/worklog" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"오늘 회의 및 개발 업무"}' | jq .
echo -e "\n"

echo "5) 오늘 업무 로그 수정 (Update WorkLog)"
curl -s -X PUT "$API_URL/worklog" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"오늘 회의, API 문서 업데이트 완료"}'
echo " → 204 No Content expected"
echo -e "\n"

echo "6) 모든 업무 로그 조회 (Get All WorkLogs)"
curl -s "$API_URL/worklog" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo -e "\n"

echo "✅ 모든 테스트 완료!"
