// Copyright 2017 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"os"
	"testing"

	providercfg "github.com/oracle/oci-cloud-controller-manager/pkg/cloudprovider/providers/oci/config"
	"go.uber.org/zap"
)

func TestBuildRateLimiterWithConfig(t *testing.T) {
	qpsRead := float32(6.0)
	bucketRead := 20
	qpsWrite := float32(8.0)
	bucketWrite := 20

	rateLimiterConfig := &providercfg.RateLimiterConfig{
		RateLimitQPSRead:     qpsRead,
		RateLimitBucketRead:  bucketRead,
		RateLimitQPSWrite:    qpsWrite,
		RateLimitBucketWrite: bucketWrite,
	}

	rateLimiter := NewRateLimiter(zap.S(), rateLimiterConfig)

	if rateLimiter.Reader.QPS() != qpsRead {
		t.Errorf("unexpected QPS (read) value: expected %f but found %f", qpsRead, rateLimiter.Reader.QPS())
	}

	if rateLimiter.Writer.QPS() != qpsWrite {
		t.Errorf("unexpected QPS (write) value: expected %f but found %f", qpsWrite, rateLimiter.Writer.QPS())
	}
}

func TestBuildRateLimiterWithDefaults(t *testing.T) {
	rateLimiterConfig := &providercfg.RateLimiterConfig{}

	rateLimiter := NewRateLimiter(zap.S(), rateLimiterConfig)

	if rateLimiter.Reader.QPS() != rateLimitQPSDefault {
		t.Errorf("unexpected QPS (read) value: expected %f but found %f", rateLimitQPSDefault, rateLimiter.Reader.QPS())
	}

	if rateLimiter.Writer.QPS() != rateLimitQPSDefault {
		t.Errorf("unexpected QPS (write) value: expected %f but found %f", rateLimitQPSDefault, rateLimiter.Writer.QPS())
	}
}

func TestDisableRateLimiterConfig(t *testing.T) {
	rateLimiterConfig := &providercfg.RateLimiterConfig{
		DisableRateLimiter: true,
	}

	rateLimiter := NewRateLimiter(zap.S(), rateLimiterConfig)
	for i := 0; i < 21; i++ {
		rateLimiter.Reader.TryAccept()
		rateLimiter.Writer.TryAccept()
	}

	if !rateLimiter.Reader.TryAccept() || !rateLimiter.Writer.TryAccept() {
		t.Errorf("RateLimiter Should be disabled")
	}
}

func TestEnableRateLimiterConfig(t *testing.T) {

	rateLimiterConfig := &providercfg.RateLimiterConfig{}
	rateLimiter := NewRateLimiter(zap.S(), rateLimiterConfig)

	for i := 0; i < 21; i++ {
		rateLimiter.Reader.TryAccept()
		rateLimiter.Writer.TryAccept()
	}
	if rateLimiter.Reader.TryAccept() || rateLimiter.Writer.TryAccept() {
		t.Errorf("RateLimiter Should be enabled")
	}

}

func TestIsIpv6SingleStackCluster(t *testing.T) {
	// Set up test cases
	tests := []struct {
		name         string
		envValue     string
		shouldSetEnv bool
		want         bool
	}{
		{
			name:         "Returns true when cluster IP family is IPv6",
			envValue:     "ipv6",
			shouldSetEnv: true,
			want:         true,
		},
		{
			name:         "Returns false when cluster IP family is not set",
			shouldSetEnv: false,
			want:         false,
		},
		{
			name:         "Returns false when cluster IP family is not IPv6",
			envValue:     "ipv4",
			shouldSetEnv: true,
			want:         false,
		},
		{
			name:         "Returns false when cluster IP family is mixed case",
			envValue:     "IPv4",
			shouldSetEnv: true,
			want:         false,
		},
		{
			name:         "Returns true when cluster IP family is mixed case IPv6",
			envValue:     "IPv6",
			shouldSetEnv: true,
			want:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up any existing env variable after test
			defer os.Unsetenv(ClusterIpFamilyEnv)

			// Set environment variable if needed
			if tt.shouldSetEnv {
				os.Setenv(ClusterIpFamilyEnv, tt.envValue)
			} else {
				os.Unsetenv(ClusterIpFamilyEnv)
			}

			// Run the function
			got := IsIpv6SingleStackCluster()

			// Validate the result
			if got != tt.want {
				t.Errorf("IsIpv6SingleStackCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
