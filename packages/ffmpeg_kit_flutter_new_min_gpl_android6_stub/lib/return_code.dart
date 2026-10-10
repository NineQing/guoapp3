class ReturnCode {
  const ReturnCode(this.value);

  final int value;

  int getValue() => value;

  static bool isSuccess(ReturnCode? code) => code?.value == 0;

  static bool isCancel(ReturnCode? code) => code?.value == 255;
}
