// Copyright 2022 Mineiros GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hcl_test

import (
	"testing"

	"github.com/mineiros-io/terramate/errors"
	"github.com/mineiros-io/terramate/hcl"
)

func TestHCLParserVendor(t *testing.T) {
	for _, tc := range []testcase{
		{
			name: "empty vendor",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						}
					`,
				},
			},
			want: want{
				config: hcl.Config{
					Vendor: &hcl.VendorConfig{},
				},
			},
		},
		{
			name: "vendor.dir",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  dir = "/dir"
						}
					`,
				},
			},
			want: want{
				config: hcl.Config{
					Vendor: &hcl.VendorConfig{
						Dir: "/dir",
					},
				},
			},
		},
		{
			name: "empty manifest",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						  }
						}
					`,
				},
			},
			want: want{
				config: hcl.Config{
					Vendor: &hcl.VendorConfig{
						Manifest: &hcl.ManifestConfig{},
					},
				},
			},
		},
		{
			name: "empty manifest.default",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						    }
						  }
						}
					`,
				},
			},
			want: want{
				config: hcl.Config{
					Vendor: &hcl.VendorConfig{
						Manifest: &hcl.ManifestConfig{
							Default: &hcl.ManifestDesc{},
						},
					},
				},
			},
		},
		{
			name: "default has files",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      files = ["/", "/test"]
						    }
						  }
						}
					`,
				},
			},
			want: want{
				config: hcl.Config{
					Vendor: &hcl.VendorConfig{
						Manifest: &hcl.ManifestConfig{
							Default: &hcl.ManifestDesc{
								Files: []string{"/", "/test"},
							},
						},
					},
				},
			},
		},
		{
			name: "files is not list fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      files = "not list"
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 13, 67), end(5, 18, 72)),
					),
				},
			},
		},
		{
			name: "vendor.dir is not string fails",
			input: []cfgfile{
				{
					filename: "vendor.tm",
					body: `
						vendor {
						  dir = ["/dir"]
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("vendor.tm", start(3, 9, 24), end(3, 12, 27)),
					),
				},
			},
		},
		{
			name: "vendor.dir is undefined fails",
			input: []cfgfile{
				{
					filename: "vendor.tm",
					body: `
						vendor {
						  dir = undefined
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("vendor.tm", start(3, 9, 24), end(3, 12, 27)),
					),
				},
			},
		},
		{
			name: "redefined vendor fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						}
						vendor {
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(4, 7, 30), end(4, 15, 38)),
					),
				},
			},
		},
		{
			name: "redefined manifest fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						  }
						  manifest {
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 9, 53), end(5, 19, 63)),
					),
				},
			},
		},
		{
			name: "redefined default fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						    }
						    default {
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(6, 11, 77), end(6, 20, 86)),
					),
				},
			},
		},
		{
			name: "files is not string list fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      files = ["ok", 666]
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 13, 67), end(5, 18, 72)),
					),
				},
			},
		},
		{
			name: "files has undefined element fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      files = ["ok", ns.undefined]
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 28, 82), end(5, 30, 84)),
					),
				},
			},
		},
		{
			name: "files is undefined fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      files = local.files
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 21, 75), end(5, 26, 80)),
					),
				},
			},
		},
		{
			name: "unrecognized attribute on vendor fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						    unknown = true
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(3, 11, 26), end(3, 18, 33)),
					),
				},
			},
		},
		{
			name: "unrecognized attribute on manifest fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    unknown = true
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(4, 11, 45), end(4, 18, 52)),
					),
				},
			},
		},
		{
			name: "label on vendor fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor "label" {
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(2, 14, 14), end(2, 21, 21)),
					),
				},
			},
		},
		{
			name: "label on manifest fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest "label" {
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(3, 18, 33), end(3, 25, 40)),
					),
				},
			},
		},
		{
			name: "unrecognized block on vendor fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  unknown {
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(3, 9, 24), end(3, 18, 33)),
					),
				},
			},
		},
		{
			name: "unrecognized block on manifest fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    unknown {
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(4, 11, 45), end(4, 20, 54)),
					),
				},
			},
		},
		{
			name: "unrecognized attribute on default fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      unknown = true
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 13, 67), end(5, 20, 74)),
					),
				},
			},
		},
		{
			name: "unrecognized block on default fails",
			input: []cfgfile{
				{
					filename: "manifest.tm",
					body: `
						vendor {
						  manifest {
						    default {
						      unknown {
						      }
						    }
						  }
						}
					`,
				},
			},
			want: want{
				errs: []error{
					errors.E(hcl.ErrTerramateSchema,
						mkrange("manifest.tm", start(5, 13, 67), end(5, 22, 76)),
					),
				},
			},
		},
	} {
		testParser(t, tc)
	}
}
