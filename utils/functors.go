package utils


func Map(values []interface{}, f func(interface{}) interface{}) []interface{}{
  mapped := make([]interface{}, len(values))

  for k, v := range values {
    mapped[k] = f(v)
  }

  return mapped
}

func Sum(v float64, v2 float64) float64 {
  return v+v2
}

func ReduceNumbers(values []interface{}, f func(float64, float64) float64) float64{
  var acc float64 = 0

  for k, v := range values {
    if k+1 < len(values) {
      acc  = acc + f(v.(float64), values[k+1].(float64))
    } else if k == 0 {
      acc = v.(float64)
    }
  }

  return acc
}
