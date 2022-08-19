# Interactive-Console

### 简介

这是一个实验程序, 用于创建一个交互式终端框架, 功能大纲如下:

- 输入输出分离, 给予用户一个分离的输入行, 避免输入输出混乱问题
- 对输入完全可控, 允许在任何时候阻止用户输入
- 命令式输入, 每个输入回车后将如命令一样产生回调

### 模拟用户接口

```cpp
void handle(const std::string &cmd, InteractiveConsole *console) {
    // ... 回调逻辑, 允许在此执行console的runFunc和write方法
    if(cmd == "stop") {
        console.stop();
    } else {
        console.Write("You typed" + cmd);
    }
}

int main() {
    // 这里有单例设计模式, 一个程序只有一个console
    InterativeConsole *console = new InterativeConsole(&handle);

    // 接收信息, 回调Handle, 阻塞知道console被关闭
    console.run();

    return 0;
}

```