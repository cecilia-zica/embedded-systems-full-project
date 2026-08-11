#pragma once

// Classification of a single reading.
struct Classification {
  int cls;      // 0=normal, 1=alert, 2=error
  float confidence; // 0.0 to 1.0
};
