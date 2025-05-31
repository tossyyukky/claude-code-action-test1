#!/bin/sh


# Claude Code で Amazon Bedrock を使用するための環境変数を設定
export CLAUDE_CODE_USE_BEDROCK=1

# 使用するモデルの指定
export ANTHROPIC_MODEL='apac.anthropic.claude-sonnet-4-20250514-v1:0'

# AWS 認証情報の設定
export AWS_REGION='ap-northeast-1'
export AWS_PROFILE='claude-code'

claude
