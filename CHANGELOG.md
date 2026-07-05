# 更新日志

本项目遵循清晰、面向用户的更新记录。每次发布前必须更新本文件，并确保版本号与 `researcher/VERSION` 和 Git tag 对齐。

## [1.0.1] - 2026-07-05

### 修复

- 默认配置读取现在固定使用 `~/.config/researcher/config.yaml`，不再受 `XDG_CONFIG_HOME` 影响。
- 更新配置读取说明，避免用户继续按旧路径排查。

### 验证

- 增加回归测试，确认设置 `XDG_CONFIG_HOME` 时仍会使用 home 目录下的 researcher 配置。
- 通过真实 Bocha 搜索和 Volcengine 问答闭环验证配置可正常加载并输出结果。

## [1.0.0] - 2026-05-28

### 新增

- 初始发布行业研究引擎。
- 提供 `researcher` Go 研究引擎，用于生成研究工作区、命题图、证据台账、反证记录、置信度报告和最终报告。
- 提供智能体技能入口，用于编排行业研究、多角色审查和证据校验流程。
- 支持 Bocha 直接网页搜索和 Volcengine Ark 联网回答检索能力。

### 发布

- 约定使用 `v1.0.0` tag 发布首个 GitHub Release。
- 发布产物包含 macOS、Linux、Windows 的 `researcher` 二进制压缩包和 SHA256 校验文件。
