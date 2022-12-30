;;;
;;; Copyright (c) 2022 Markku Rossi
;;;
;;; All rights reserved.
;;;

(define (test name)
  (letrec ((count 0)
           (success 0)
           (fail 0)
           (run-test
            (lambda (category idx num-tests test)
              (display " - test ") (display category) (display ": ")
              (display (+ idx 1)) (display " / ") (display num-tests)
              (display " ")
              (let ((failed #f))
                ;; Run test with reporter.
                (test (lambda (cmd . args)
                        (cond ((eq? cmd 'error)
                               (set! failed #t))
                              (else
                               (display "Test reporter: invalid command: ")
                               (display cmd)
                               (newline)))))
                (if failed
                    (begin
                      (set! fail (+ fail 1))
                      (display #\x274c))
                    (begin
                      (set! success (+ success 1))
                      (display #\x2713)))
                (newline))))
           (iter
            (lambda (category num-tests tests)
              (if (null? tests)
                  #t
                  (begin
                    (set! count (+ count 1))
                    ;;((car tests) runner)
                    (run-test category
                              (- num-tests (length tests))
                              num-tests
                              (car tests))
                    (iter category num-tests (cdr tests))))))
           (runner
            (lambda (cmd . args)
              (cond ((eq? cmd 'run)
                     (car args)
                     (iter (car args) (length (cdr args)) (cdr args)))
                    ((eq? cmd 'error)
                     (set! fail (+ fail 1)))
                    ((eq? cmd 'stats)
                     (display "Tests  : ") (display count) (newline)
                     (display "Success: ") (display success) (newline)
                     (display "Fail   : ") (display fail) (newline))
                    ((eq? cmd 'status)
                     (= fail 0))
                    (else
                     (display "Usage:") (newline)
                     (display " - run tests...") (newline)
                     (display " - stats") (newline))))))
    runner))

(define runner (test "r6rs"))

(load "test-characters.scm")
(load "test-eq.scm")

(runner 'stats)
(runner 'status)
