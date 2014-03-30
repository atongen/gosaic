package command

// hasExpectedArgs checks whether the number of args are as expected.
func hasExpectedArgs(args []string, expected int) bool {
  switch expected {
  case -1:
    if len(args) > 0 {
      return true
    } else {
      return false
    }
  default:
    if len(args) == expected {
      return true
    } else {
      return false
    }
  }
}
