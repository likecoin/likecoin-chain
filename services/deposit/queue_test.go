package deposit

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueue(t *testing.T) {
	Convey("For a newly initialized queue with capacity 16", t, func() {
		q := newQueue(16)
		Convey("Its size should be zero", func() {
			So(q.size(), ShouldEqual, 0)
			Convey("Peek should fail", func() {
				_, ok := q.peek()
				So(ok, ShouldBeFalse)
			})
		})
		Convey("After enqueue", func() {
			q.enqueue(1337)
			Convey("Its size should be 1", func() {
				So(q.size(), ShouldEqual, 1)
				Convey("Peek should return the value", func() {
					n, ok := q.peek()
					So(ok, ShouldBeTrue)
					So(n, ShouldEqual, 1337)
					Convey("After some further enqueue", func() {
						for i := uint64(1); i <= 100; i++ {
							q.enqueue(i)
						}
						Convey("Its size should equal to the number of enqueue called", func() {
							So(q.size(), ShouldEqual, 101)
							Convey("Peek should return the first enqueued value", func() {
								n, ok := q.peek()
								So(ok, ShouldBeTrue)
								So(n, ShouldEqual, 1337)
								Convey("After dequeue, peek should return by the order of enqueue", func() {
									for i := uint64(1); i <= 100; i++ {
										q.dequeue()
										n, ok := q.peek()
										So(ok, ShouldBeTrue)
										So(n, ShouldEqual, i)
										So(q.size(), ShouldEqual, 101-i)
									}
								})
							})
						})
					})
				})
			})
		})
		for i := 0; i < 100; i++ {
			Convey(fmt.Sprintf("After some random operations (iteration %d)", i), func() {
				ops := uint64(rand.Intn(9000) + 1000)
				ch := make(chan uint64, 10000)
				for j := uint64(0); j < ops; j++ {
					switch rand.Intn(2) {
					case 0:
						q.enqueue(j)
						ch <- j
					case 1:
						q.dequeue()
						if len(ch) > 0 {
							<-ch
						}
					}
				}
				Convey("The size of the queue should be correct", func() {
					size := q.size()
					So(size, ShouldEqual, len(ch))
					Convey("The contents should be correct", func() {
						for j := 0; j < size; j++ {
							n, ok := q.peek()
							So(ok, ShouldBeTrue)
							So(n, ShouldEqual, <-ch)
							q.dequeue()
							So(q.size(), ShouldEqual, len(ch))
						}
						_, ok := q.peek()
						So(ok, ShouldBeFalse)
						So(q.size(), ShouldEqual, 0)
						q.dequeue()
						So(q.size(), ShouldEqual, 0)
					})
				})
			})
		}
	})
}
