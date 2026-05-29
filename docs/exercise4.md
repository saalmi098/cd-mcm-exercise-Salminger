# CDMA Exercise 4: SALMINGER

## Task 1: Docker Image Scanning

### Before (alpine:3.19)

Total: 10 vulnerabilities — CRITICAL: 0, HIGH: 2, MEDIUM: 5, LOW: 3

![](img/ex4-task1-1.jpg)

All 10 from `alpine:3.19` base (busybox, musl, ssl_client). Go binary: 0 vulns.

### Optimization: Switch to distroless

Changed runtime base to `gcr.io/distroless/static-debian12` — no shell, no package manager, no busybox/musl.

### After (distroless/static-debian12)

Total: 0 vulnerabilities

![](img/ex4-task1-2.png)

### CI: Trivy Scan Job

Trivy report aftifact: [https://github.com/saalmi098/cd-mcm-exercise-Salminger/actions/runs/26636398636/artifacts/7290317405](https://github.com/saalmi098/cd-mcm-exercise-Salminger/actions/runs/26636398636/artifacts/7290317405)
