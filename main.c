#include "9cc.h"

void error(char *fmt, ...) {
  va_list ap;
  va_start(ap, fmt);
  vfprintf(stderr, fmt, ap);
  fprintf(stderr, "\n");
  exit(1);
}

char *user_input;
Token *token;
Node *node;

int main(int argc, char **argv) {
  if (argc != 2) {
    error("引数の個数が正しくありません");
    return 1;
  }

  // トークナイズしてパースする
  user_input = argv[1];
  token = tokenize();
  Node *node = expr();

  // アセンブリの前半部分を出力
  printf(".intel_syntax noprefix\n");
  printf(".globl main\n");
  printf("main:\n");

  // 抽象構文木を下りながらコード生成
  gen(node);

  // スタックの一番上にある値が答えなので、popして返す
  printf("  pop rax\n");
  printf("  ret\n");
  
  // NOTE: GAS (GNU Assembler)にスタックを実行不可にさせる
  // NOTE: 現状スタックを使用する必要がないのと、無効化しないと警告が発生するため追加
  printf(".section .note.GNU-stack,\"\",@progbits\n");
  return 0;
}
