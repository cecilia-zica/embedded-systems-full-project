#pragma once

// se você já incluiu esse arquivo antes (nesse mesmo arquivo final), não inclua de novo
struct Classificacao {
  int classe;      // 0=Normal, 1=Alerta, 2=Erro
  float confianca; // 0.0 a 1.0
};
