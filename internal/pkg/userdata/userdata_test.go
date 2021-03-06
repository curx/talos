/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package userdata

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// nolint: lll
const testConfig = `version: ""
security:
  os:
    ca:
      crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNMRENDQVk2Z0F3SUJBZ0lRTHhqVGhwbjB0SGpzZ25xMGlJV21HakFLQmdncWhrak9QUVFEQkRBak1RNHcKREFZRFZRUUtFd1ZVWVd4dmN6RVJNQThHQTFVRUF4TUlWR0ZzYjNNZ1EwRXdIaGNOTVRneE1qQTFNVFF6TWpFeApXaGNOTVRreE1qQTFNVFF6TWpFeFdqQWpNUTR3REFZRFZRUUtFd1ZVWVd4dmN6RVJNQThHQTFVRUF4TUlWR0ZzCmIzTWdRMEV3Z1pzd0VBWUhLb1pJemowQ0FRWUZLNEVFQUNNRGdZWUFCQUFRN0I4TWhxZG1IZDduOXVKT3RLTU8Kb0dOamJ0U3YxcHk1QUdPMTQzd3hsVkQ3RWpKK215KzZGQkk1T3ZEakdSYnFZbnNiWFdrWGorWG5qZmpPcnBYNwo2d0UwbEF6aEJwMzRnWENsUC92Z2FzWlk2dmJlYnUwSGdTQ2tLYW9IMlgrNDlBNWEvNHBhTURsSjBCSGlQa3plCmVIcHV0TkVHWitTYThEcHpFaitxdnp4SkxLTmhNRjh3RGdZRFZSMFBBUUgvQkFRREFnS2tNQjBHQTFVZEpRUVcKTUJRR0NDc0dBUVVGQndNQkJnZ3JCZ0VGQlFjREFqQVBCZ05WSFJNQkFmOEVCVEFEQVFIL01CMEdBMVVkRGdRVwpCQlJwRDFubnBrSWJWSldYUC9YZDRIR0xRSEh2UERBS0JnZ3Foa2pPUFFRREJBT0Jpd0F3Z1ljQ1FVLzNkbm1ECi8wNmJXaHNVcE0zM2hUWnhwM3NTcDlKTHF6Zm5xUUMrTTBQYWwrMHc1ZzEyb0RuN0s3RnVZcnM1UGFlcy9MREIKa0N1Y3VvZmxESG82OFRrVUFrSUFob01MMXhXUVU3Z1UvODMyY1JuV3A2V0xnK3VWRGVWQno3YW1UK21rT1RYYgpqS2dpRCs5bmlEc3U3elFoS081NVFKZTFEU1JsQzQ1cTRJZ01NQzFnckE4PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
      key: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1JSGNBZ0VCQkVJQjY1cmZJUmZ6NVRPR0VuWVVnL05iSWU0NFhVelFCNG9DWU9qYWY2WU83TktsSS9WTmE0aDYKUFJlY0Vsay9LZzZTK1pXWTB3MGwxSC9JTVFXNXhzSDVJQnlnQndZRks0RUVBQ09oZ1lrRGdZWUFCQUFRN0I4TQpocWRtSGQ3bjl1Sk90S01Pb0dOamJ0U3YxcHk1QUdPMTQzd3hsVkQ3RWpKK215KzZGQkk1T3ZEakdSYnFZbnNiClhXa1hqK1huamZqT3JwWDc2d0UwbEF6aEJwMzRnWENsUC92Z2FzWlk2dmJlYnUwSGdTQ2tLYW9IMlgrNDlBNWEKLzRwYU1EbEowQkhpUGt6ZWVIcHV0TkVHWitTYThEcHpFaitxdnp4SkxBPT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=
    identity:
      crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNZRENDQWNHZ0F3SUJBZ0lSQUl5cjIwVXNhaXpTOWZtK1FlOFZMSVV3Q2dZSUtvWkl6ajBFQXdRd0l6RU8KTUF3R0ExVUVDaE1GVkdGc2IzTXhFVEFQQmdOVkJBTVRDRlJoYkc5eklFTkJNQjRYRFRFNE1USXdOVEUwTXpVMQpObG9YRFRFNU1USXdOVEUwTXpVMU5sb3dKVEVPTUF3R0ExVUVDaE1GVkdGc2IzTXhFekFSQmdOVkJBTVRDbFJoCmJHOXpJRTV2WkdVd2dac3dFQVlIS29aSXpqMENBUVlGSzRFRUFDTURnWVlBQkFHTHl1aXFIcWpzdWhMcE1sZWcKL2pOWEdXYTdBN3pIRWRvSG9KTnNWcmRPTHFBUldlRjNULzkrRVM2aXRISXRBWU5ONis0THFITHQ0dkdmT1E0eAp3OUdpR2dEdmwrUlhQZzVibUVNby90ZmxRNHk0QVdTK3R2Z2F2MG9hZXlOSmh0YVk2Y0hrQ3RVR2tXcGRKMEVnCllFQ2xQcFhnbXRSMVVMSGQ1dTdKeGQ0NFp2R1hTNk9Ca0RDQmpUQU9CZ05WSFE4QkFmOEVCQU1DQmFBd0hRWUQKVlIwbEJCWXdGQVlJS3dZQkJRVUhBd0VHQ0NzR0FRVUZCd01DTUF3R0ExVWRFd0VCL3dRQ01BQXdId1lEVlIwagpCQmd3Rm9BVWFROVo1NlpDRzFTVmx6LzEzZUJ4aTBCeDd6d3dMUVlEVlIwUkJDWXdKSWNFQ21PQ2pJY0VDbU9DCmpZY0VDbU9Dam9jRWswdE8yWWNFazB0THZZY0VrMHRMaVRBS0JnZ3Foa2pPUFFRREJBT0JqQUF3Z1lnQ1FnRzcKdnU2MTBNcktLdytFMVBTaEtLMkJkZVV4RmRvYXRBV1U2L3Z6RjZ5MDZrci9XL0xvc04remtuSWRsT1FoZDVIeQpyMWh3Mm1tQXRtVjNKYTNkVFNYUExBSkNBZFYzRUlWaUtvb2JnTStySTl4S1dPbmRMS2hLK3VoUHdXZy9QNmx4ClNjcU1sMlo0WFhjZ3NudnppdDZBZFBMV3ZmUkN2WjF5NW9LTzVFc0IvOGU0Um1nbgotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
      key: LS0tLS1CRUdJTiBFQyBQUklWQVRFIEtFWS0tLS0tCk1JSGNBZ0VCQkVJQWc0akVSbUhTM2RHajIyYVcwRlVFS0ZBNWEwazk5cHBmSjNnb2FTQll1aS9mRHZXb0YveUYKQnZUeHh4L0w2UFZhSmZ0UVFaZUVuTXVrU3IrL1FacnhPYTZnQndZRks0RUVBQ09oZ1lrRGdZWUFCQUdMeXVpcQpIcWpzdWhMcE1sZWcvak5YR1dhN0E3ekhFZG9Ib0pOc1ZyZE9McUFSV2VGM1QvOStFUzZpdEhJdEFZTk42KzRMCnFITHQ0dkdmT1E0eHc5R2lHZ0R2bCtSWFBnNWJtRU1vL3RmbFE0eTRBV1MrdHZnYXYwb2FleU5KaHRhWTZjSGsKQ3RVR2tXcGRKMEVnWUVDbFBwWGdtdFIxVUxIZDV1N0p4ZDQ0WnZHWFN3PT0KLS0tLS1FTkQgRUMgUFJJVkFURSBLRVktLS0tLQo=
  kubernetes:
    ca:
      crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUZSVENDQXkyZ0F3SUJBZ0lRV3dqZWpSUldmc2VEaFl6UkRwcXlCekFOQmdrcWhraUc5dzBCQVFzRkFEQXQKTVJNd0VRWURWUVFLRXdwTGRXSmxjbTVsZEdWek1SWXdGQVlEVlFRREV3MUxkV0psY201bGRHVnpJRU5CTUI0WApEVEU0TVRJd05URTBNekl4TTFvWERURTVNVEl3TlRFME16SXhNMW93TFRFVE1CRUdBMVVFQ2hNS1MzVmlaWEp1ClpYUmxjekVXTUJRR0ExVUVBeE1OUzNWaVpYSnVaWFJsY3lCRFFUQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQUQKZ2dJUEFEQ0NBZ29DZ2dJQkFLWlh4aFRlNnc1MkRYS2NMSVR2VGM3ZTc2eWhpcG5Zb0tndVdtNWMzT0xqTVpkawpPTkxra0xKK01QNzR3QlpRaGlmaWFCZ0xhY0FobXMzVHZrZEMzRkFNSnpiRXJmd25nazhNQnBPMmtBRW9FNWFpCmxhNEJTYWlIZXcvWkZUamE2elljRXVLdExDa2ZwSEFvVU94NUtYQ2hGY1RmVUh4U1ZndUJXaWJYMERHNjltTkwKT1dhQVBqSEw1NWFsWC9TcDF0NHU4OXNIaDZmZVRFc2o0aCtSMlBVZUVkS3I3NHdRM21EWXB4Rkl1NzFvNUxOTgprKzNYRDNLZldJZGlYUTNnV1UwSXU4czVFYmkrVUg1OC9jRW1GUUpWNm8wQ1lhcnkzRVRGWEtYaS9RUXVrUEhOCjh3YkZlVWJaYktIcm5VU2JnNllpNDRVUTg2RUJxZm90b2JmSHZpZmsvUlY1amFSQko0cjhTakRTeFJyd1A2VXkKNlZISitaS29XVFB5M1kxYTFoL3VFTEhjaDliNFdJTUtWaW9BOTFrTjFmcytzMkN3NXJzV2x5eHpXRWpCdUpGLwpPTDVRR2ZLbW1sREdGT3N6MXVZSjhNcmxwVkp1b0NiQU8ra3JWWEhwWnhtcUpPS3JwWEVicUg1UHMwbUlCZkYzCk5wUElJMEFXbUd2bnBudzQxaDBUc0dmcHFjak9ZclM3TXE1RkYrSTNHbUM5OGV4Q1loYnFPNGVsK01WNk9PY3oKQ2dRenBTdDlTYUV4NHFKaG5TaWthbGNiR3BHYlFtNERyMFdvcnkxOTI2Wjl0WlhCRS91ZDJQR0Q2OFBKSi9NbgozUHhHbmNoUVdhbVRiaWpEay9QM3F4OWhuM0d1MEVaWSsvMkEzSGpRTGxKTkhtNzlhb1YyQkdEazRPNjNBZ01CCkFBR2pZVEJmTUE0R0ExVWREd0VCL3dRRUF3SUNwREFkQmdOVkhTVUVGakFVQmdnckJnRUZCUWNEQVFZSUt3WUIKQlFVSEF3SXdEd1lEVlIwVEFRSC9CQVV3QXdFQi96QWRCZ05WSFE0RUZnUVVQcW1iVk9oL3VkbkRscWxlVS9zVgpncTI0V2xRd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dJQkFKL0JYWU9JM2hJemxUSTlaSFlSclpGenlRR0wwbExaCklCYnB0UVhWN1Vla2x1RUJBbkRCU21VUytETmh6ZTJVZGlzUHdaV1oxZWtYNnRiT2ZwYVJGWlYwT1IxTm5UU2IKQmFQNUdyeW9JYXVhcjhVT1BkRkI4b0RsYWQrOHQwKzNZVFQvRyt5dC9tdlVReEJ1c3g3ckJLbE1nS1F3ejBqdQpEUi84Y1BuY0JSa2h6eExqRnJFY2REb1BDNi83MmJUTzF2S0VOczZOdVBtb1h0TzJVdkhUdlBrUWZVR0dyYzRVCjZEZjhjTTZBSnlDRnNIM0NUZUdxRGsvK3ZjQWFpcUZRNW9rZTRvOGxqYklqbkFGTytXWmhMMkg5UGMwRFlybEEKZVRUbUNxNitCa2dkRmJFZmpXdGY1NDJXdlJGRjZPSm9qcEsveUlZNlRjQTlSa0o1blY5UHFCcTVEdG1oeW9uNwpzalN0TWdxRWVUdDNVUC9ZMDI1eEN3T1orMU9halpMdjUwTHhkbVFQNkh4NkJrYzFscDZSSENxTTdQZXhhZlJSClZRODFmVVhURzl5dkNxTmdpdXdWYzloMkRsb25PUGdMaGNDS0ozYjE2NWE1TXFMMkdRZ3BOY254dC9wblE5SGUKeVFiekMwZFhDOXlKSWc2MTFWZ3lkMkJWVlFISWsvTU1lRmlEM2RWa1ZKSW95WU9tcy9nQmN4dHM2TzkwRnY5QwpxNUI2TUIvc016a2pLYnpXNFVPb2ltd2grVUtVRVVaUVhXbGs3UXN3U3NmT3Fidyt4ajRxUzl4MmFvRXdBSXM2Ckloa0NZZmRqZU1henFmV21valhQM1FNYTZsNXZuRlZMa29ZOTRIT285RDdjQi9pSGlMbXEzejRESGlYamU3RWQKVENlWG1NZ25CcnkzCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
      key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS1FJQkFBS0NBZ0VBcGxmR0ZON3JEbllOY3B3c2hPOU56dDd2cktHS21kaWdxQzVhYmx6YzR1TXhsMlE0CjB1U1FzbjR3L3ZqQUZsQ0dKK0pvR0F0cHdDR2F6ZE8rUjBMY1VBd25Oc1N0L0NlQ1R3d0drN2FRQVNnVGxxS1YKcmdGSnFJZDdEOWtWT05yck5od1M0cTBzS1Ira2NDaFE3SGtwY0tFVnhOOVFmRkpXQzRGYUp0ZlFNYnIyWTBzNQpab0ErTWN2bmxxVmY5S25XM2k3ejJ3ZUhwOTVNU3lQaUg1SFk5UjRSMHF2dmpCRGVZTmluRVVpN3ZXamtzMDJUCjdkY1BjcDlZaDJKZERlQlpUUWk3eXprUnVMNVFmbno5d1NZVkFsWHFqUUpocXZMY1JNVmNwZUw5QkM2UThjM3oKQnNWNVJ0bHNvZXVkUkp1RHBpTGpoUkR6b1FHcCtpMmh0OGUrSitUOUZYbU5wRUVuaXZ4S01OTEZHdkEvcFRMcApVY241a3FoWk0vTGRqVnJXSCs0UXNkeUgxdmhZZ3dwV0tnRDNXUTNWK3o2ellMRG11eGFYTEhOWVNNRzRrWDg0CnZsQVo4cWFhVU1ZVTZ6UFc1Z253eXVXbFVtNmdKc0E3NlN0VmNlbG5HYW9rNHF1bGNSdW9mayt6U1lnRjhYYzIKazhnalFCYVlhK2VtZkRqV0hST3daK21weU01aXRMc3lya1VYNGpjYVlMM3g3RUppRnVvN2g2WDR4WG80NXpNSwpCRE9sSzMxSm9USGlvbUdkS0tScVZ4c2FrWnRDYmdPdlJhaXZMWDNicG4yMWxjRVQrNTNZOFlQcnc4a244eWZjCi9FYWR5RkJacVpOdUtNT1Q4L2VySDJHZmNhN1FSbGo3L1lEY2VOQXVVazBlYnYxcWhYWUVZT1RnN3JjQ0F3RUEKQVFLQ0FnQWs0Y05PcjFxSTIwNENBblN3aU9yRW1wT2p3REdlQ1BVZU5TRGg1WDhvWTEyRWhyaytzV1VQM29ENApsNmpuaWJVbE5NTUZ5Y29KeXFtclIyNmlHRVlIRFpySVB2V2d1aFhmZHZnaVdsRTFvSEF2Yng4bTlLd2pUTjdSClZiMnAzSWhZZUFNTDlYK2NJcUx3SjFCQ0RsOU0xTHFoNnkzS1V1czNJOWdjVXErYnh5dUlPbzZnbG9Denc5VTIKaHJadTVoVjVNQW5ybjVESmZMV0gxNDMwbC83MTRsUHJWU29JOFZpTk9weFliYnNLdFM0NFAxUTBZUUhRMVZubwpzcnVWTHhBOXoxanNKMXd4alNrRjBxcmszNCs2ODlmTXFpR0RxTk9FQUFxWjBXbHVPdkR5WEVxdFRxMmtNdU11ClI3S1lHcThtOEFkb1lxalNJeTlRQkR3ZUpwQ1hwa1pxS1Btc1ozZVlGMXBOMHFFbWFaUUxFNHNDdjd5anpSd3oKSXpPS21DOS9UcFVOaE94Z0E1UkExUEpoWVZZYlVwc0k2bjEzOEVqSjY1cHpSTW55V01VZlZIQ0xPVGdUTFdmeApTRWNKOW1vZW00emwvclRxZWJmbm9WMm54elpQYTlCMnFmMlFVYXpBKzZxTDY4U2l5c2l5bXB4d2pzMWVjTTM2ClkzMDZUMFBrUkw1MFhiTDYvTk5uRGxpMmRrMlVHTzlUVTFMeDZsTlAvd2hPbzFqN1AzUm5OelFzRDB3aG5MakQKL05zUmFpRndzYVhKUU4rUzBRVEN1cTZWU1NYTjZaK3FmMWNsdkh3aXdvZERZbGRsR3lYUWN2VWErMmg3M3Jadgo2TW1nRm13Znc3OWhrS1B4UjdaOXdRL0JLeW1PMEVjS3FyWnlvYWpqUC9lZHRHMGFhUUtDQVFFQTJGcVk0VWo0CjN0QVJXcVRGQTRwYXFUNHlzZU14M0g0dFhTWm8xQ200cEFWc1pMNFEvYkRYNFZCTUNzUFRnVVdseWtJZlJJNEEKRC9zaHYvSGFHNkZRSTUvQUlza2ZaMWhYbzBUOGRzaHVYeUJkbUc0bk56cEpxYzF4UktyS0dsYTA3a2lDVHd0NwpncUxJNldmZHVaMHhZVTJtOFFPTTNWckpRVmpmM1lUczN0eitqY3RFRjJoYi93UDhJZzNVTCtpbmJOb0NNR3FZCm9kSWY4Wlk0VUdnV3NiYlA0K29CUUtnUHljcjRxWGUvayt6ZXJPaitoQ0hrV3NpWVhiay9hMm40V3FIWGxrUDUKaVRRbGtMTzQ1UXVxc2hkTWZ5Z2NDZ3JQUWNGNXZGQ01XTTJMaHJWSTBISmx6UjJCL2FLaERqZmxNY2l2VU1JdApISHJYbTJUcnRrN2Znd0tDQVFFQXhOTWFkVnB3SWsxVjBPT3RwYlNjVXlkWlAyNnlkUGl1a2NHcXhHSVcydVR2ClpqanNOQTY5UlpLbFhsVWNmVVdmY3dob1U1blcySzg4My9lOUZUanpZK0ZrTVc4VW5hN005OEdxTlhkQ09pSzAKUHcvK0huTU1DTlpyaVloVDR0RnFNYTdWdjhMNzBBRmxBMnB3MzEyR0FiSm1hRVF6ZkRNQ3RvUVRZdnVpRFhUcworSzJTaHVyQW5yUWRmNE81dkpwYXlzRkJJMFhOalo2WlpndU0wMDdhYm9ZcW4vS0V1cVVCTUh2OFQ1Q1pxM0ZZClBFYUtZZitrSWhmWnRjc3QzYWlXUE4zZ3dDQzFrNHBlaWtPMzk0ZlRqOVdmTVhZa0c1cUNIaVUwUFo5ZHU3QkMKYk1oS21IcHhCNklCS2FpaXgwTzd6cUxreXAvSGVHY0ZqbmI2Zml0NXZRS0NBUUFOK2p4cVFaYWlmbnJBaW1pWQpBL1k3Zk9NMWp1SUh4cmNUajRteU8wZk1nUFV4eFAzQUJnN01aYTJqL0diTHNUNDJ4UExVTVFCY0IvTjBQU0hFCkt6WE1OMlBvVzJvRitUVWdQVEs2VWRTZm5LMnZUVjZIT09MTmI1Smp4MHpyU3JMQnVqbUE5ZEx5NjZWalB1eWEKTTBlZmE2N2ZYMFZZZjZjRTY4TDZ3cjJ5NEVBcDFQbi9NU1RRWXVlRkk3T1RyTW8wUkJsa0cxN2xCWGcrMlYrbQpBak9GSTdSbW14V2RvYjN6WVlPVEgwTm5RU1JaczJ2T0NZcUJPdmh0QmF6Tk9ibHIwWXptRGxvdXZRbTRRWVF1CmVBUjlJUGcyTnRjbzV1M2c1NmovdTR4MXFFSGRZQlRtTXAzVkZKVHpWL0JqeE1TdjVMRSszR3lockdZRmlnMlkKWjV3VkFvSUJBUUNKNm1wbHhNTkNKZkoxOXJPY294dk80SW1XK3hndEFueEFRZlpJdVhyT1V4QzFEam5rbUZkagp3c2dCZTU0RTRaNUMycWRQdlVKZ2Z4b216YzgrYWJPU1QwTUhxSkJIbFU2bVpKdkI3eGtHNTM1cG1lazBxT1RMCjVJTmgwWmNHSGp0TktQYzU4WVRGVWRUWU9lTnhnSFVBQjJQZXQ3dGowZEthMVoxSUJJNmU0MG1SVHY2WGJ6TWQKZ1pRUjM5bHNJTFJ4MGt5MFRKc21HdlhlSFVYR1kwQ0hwWmVFa0VrVkI3MzNuVWxFNzB2ZWFBT3Q2TGZheDJqWgpRdUVZazVCRiszV2FuTlo1aWZIeHVrbFlsTzc1a2NGYXpoTkljMElGM1pSRFQwcDF5NC9BamxkZEpsT3UyRTBQClRpT0diMVh3VHJQd3JRK29NNkZ2NnJCb0RWS2t5ZzM1QW9JQkFRQ2VFdllVeG1iRlIwTXpHTTdyZklPNWdtbm4KUStyd1d6NkJYZFJCeFdxY2FlZ09helNXYmpnaDdFS0lvVmg0THh3Ykw3a0lCaUZmVzY3a0tONUxSTnYyUzRRdwp4QnhrcHFuZ2ZoSTdhbloyZ0JkWm9NK3JEcU8wK3A2YnFiWXRKZnZrOGJUdDJQZ1pyMWg5Q0Jmbm5uaHNpUE5HCk5xTW9BN2FQTnA5cmlTS05YNDlUS096YysvREg2SkdIQUZxM0x0TVBpK2U5L01Fb3hZUFFhTjR6QnhZdE1qb2QKVWd3dVR5b2pBS3JKZmd6eDVWbmdPdzlQeEVzL0FHemxqaWQ3RHlBamR1VXBZRVNSdzJUdnh3SU1QaFEvSDJvUQo5V3l5L0FNK1BXTVU0Y2pySXBBZ2JpdFQ4ODFOR2R6c0M5bFNubXFpTVZwNkxVNkx2SklTWXZyUWNZU3AKLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K
networking:
  os: {}
  kubernetes: {}
services:
  init:
    cni: flannel
  kubeadm:
    configuration: |
      apiVersion: kubeadm.k8s.io/v1alpha3
      kind: InitConfiguration
      apiEndpoint:
        advertiseAddress: 147.75.78.217
        bindPort: 6443
      bootstrapTokens:
      - token: '1qbsj9.3oz5hsk6grdfp98b'
        ttl: 0s
      nodeRegistration:
        taints:
        kubeletExtraArgs:
          node-labels:
          feature-gates: ExperimentalCriticalPodAnnotation=true
      ---
      apiVersion: kubeadm.k8s.io/v1alpha3
      kind: ClusterConfiguration
      clusterName: talos.cluster.local
      controlPlaneEndpoint: 147.75.78.217:443
      apiServerCertSANs: [ 147.75.78.217,147.75.75.189,147.75.75.137 ]
      apiServerExtraArgs:
        runtime-config: settings.k8s.io/v1alpha1=true
        oidc-issuer-url: https://dex.dev.autonomy.io
        oidc-username-claim: email
        oidc-groups-claim: groups
        oidc-client-id: kubernetes
        feature-gates: ExperimentalCriticalPodAnnotation=true
      controllerManagerExtraArgs:
        terminated-pod-gc-threshold: '100'
        feature-gates: ExperimentalCriticalPodAnnotation=true
      schedulerExtraArgs:
        feature-gates: ExperimentalCriticalPodAnnotation=true
      networking:
        dnsDomain: talos.cluster.local
        podSubnet:
        serviceSubnet: 192.168.0.0/24
      ---
      apiVersion: kubeproxy.config.k8s.io/v1alpha1
      kind: KubeProxyConfiguration
      mode: ipvs
      ipvs:
        scheduler: lc
  trustd:
    image: docker.io/autonomy/trustd:latest
    username: '5H7iU_9\u003cxizSRD'
    password: '[a@c56X!@78\u0026f\u0026L7%+ibLrZG'
    endpoints: ["147.75.75.189"]
  proxyd:
    image: docker.io/autonomy/proxyd:latest
  blockd:
    image: docker.io/autonomy/blockd:latest
  osd:
    image: docker.io/autonomy/osd:latest
install:
  wipe: true
  boot:
    device: /dev/sda
    size: 1024000000
  root:
    device: /dev/sda
    size: 1024000000
  data:
    device: /dev/sda
    size: 1024000000
`

func TestDownloadRetry(t *testing.T) {
	ts := testUDServer()
	defer ts.Close()

	_, err := Download(ts.URL)
	if err != nil {
		t.Error("Failed to download userdata", err)
	}
}

func testUDServer() *httptest.Server {
	var count int

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
		log.Printf("Request %d\n", count)
		if count == 4 {
			// nolint: errcheck
			w.Write([]byte(testConfig))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}))

	return ts
}
