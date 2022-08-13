#define CLK 2
void setup() {
  // put your setup code here, to run once:
  pinMode(CLK, OUTPUT);  
  pinMode(LED_BUILTIN, OUTPUT);
}

bool pulse;

void loop() {
  digitalWrite(CLK, pulse);
  digitalWrite(LED_BUILTIN, pulse);
  delay(150);
  pulse = !pulse;
}
