# Copyright 2021, Staffbase GmbH and contributors.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.24.3 AS build
WORKDIR /yamllint-action
COPY go.mod go.sum /yamllint-action/
RUN go mod download
COPY . .
RUN export CGO_ENABLED=0 && go build -o ./yamllint-action .

FROM python:3.13.1-alpine3.19
RUN pip install --no-cache-dir yamllint==1.35.1 && \
    adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid 10001 appuser
COPY --from=build /yamllint-action/yamllint-action /yamllint-action
COPY entrypoint.sh /entrypoint.sh
USER appuser:appuser
ENTRYPOINT ["/entrypoint.sh"]
