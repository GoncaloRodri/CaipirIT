package utils


func Contains(slice []uint, element uint) bool {
    for _, v := range slice {
        if v == element {
            return true
        }
    }
    return false
}