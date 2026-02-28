export default {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat',     // 新功能
        'fix',      // 修复 bug
        'docs',     // 文档变更
        'style',    // 代码格式（不影响代码运行的变动）
        'refactor', // 重构（既不是新增功能，也不是修改 bug 的代码变动）
        'perf',     // 性能优化
        'test',     // 增加测试
        'chore',    // 构建过程或辅助工具的变动
        'revert',   // 回退
        'build',    // 打包
      ],
    ],
    'subject-case': [0], // 关闭主题大小写检查
    'subject-max-length': [2, 'always', 72], // 主题最大长度
    'body-max-line-length': [2, 'always', 100], // 正文每行最大长度
  },
};
