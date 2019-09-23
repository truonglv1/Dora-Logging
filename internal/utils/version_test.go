package utils

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"testing"
)

//=> done
//TODO: xem lại dùng package đặc biệt dành cho version riêng thay vi string compare tự viết
func Test_Version(t *testing.T) {
	var err error
	var v1, v2, v3, v4 *version.Version
	if v1, err = version.NewVersion("v1.1"); err != nil {
		panic(err)
	}
	if v2, err = version.NewVersion("1.2.0"); err != nil {
		panic(err)
	}

	if v3, err = version.NewVersion("1.2.3"); err != nil {
		panic(err)
	}
	if v4, err = version.NewVersion("1.3"); err != nil {
		panic(err)
	}
	fmt.Println(v1.Compare(v2))
	fmt.Println(v1.LessThan(v3))
	fmt.Println(v3.LessThan(v4))
	fmt.Println(v1)

}

func Test_BinarySearch(t *testing.T) {
	var err error
	var v, v1, v2, v3 *version.Version
	if v1, err = version.NewVersion("v1.0"); err != nil {
		panic(err)
	}
	if v2, err = version.NewVersion("v1.0.1"); err != nil {
		panic(err)
	}

	if v3, err = version.NewVersion("v2.0"); err != nil {
		panic(err)
	}

	if v, err = version.NewVersion("v1.1"); err != nil {
		panic(err)
	}

	versions := []*version.Version{}
	versions = append(versions, v3)
	versions = append(versions, v2)
	versions = append(versions, v1)

	fmt.Println(versions)
	l := 0
	r := len(versions) - 1
	for l <= r {
		mid := l + (r-l)/2
		fmt.Println("mid: ", mid)
		re := versions[mid].Compare(v)
		if re == 0 {
			fmt.Println("result: ", versions[mid])
			break
		}
		if re < 0 {
			r = mid - 1

			if mid > 0 {
				if versions[r].GreaterThan(v) {
					fmt.Println("result: ", versions[r])
					break
				}
			} else {
				if versions[0].LessThan(v) {
					fmt.Println("result: ", versions[0])
					break
				}
			}
			continue
		} else {
			l = mid + 1

			if mid < len(versions)-1 {
				if versions[l].LessThan(v) {
					fmt.Println("result: ", versions[l])
					break
				}
			}

			continue
		}
	}
}

func TestDiv(t *testing.T) {
	fmt.Println(float64(3) / 10 * 100)
}
