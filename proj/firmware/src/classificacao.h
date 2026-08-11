#pragma once

// Classification of a single reading.
struct Classificacao {
  int classe;      // 0=normal, 1=alert, 2=error
  float confianca; // 0.0 to 1.0
};
