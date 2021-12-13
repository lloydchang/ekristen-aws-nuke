package resources

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/rebuy-de/aws-nuke/pkg/types"
)

type IAMUserPolicyAttachment struct {
	svc        iamiface.IAMAPI
	policyArn  string
	policyName string
	userName   string
}

func init() {
	register("IAMUserPolicyAttachment", ListIAMUserPolicyAttachments)
}

func ListIAMUserPolicyAttachments(sess *session.Session) ([]Resource, error) {
	svc := iam.New(sess)

	resp, err := svc.ListUsers(nil)
	if err != nil {
		return nil, err
	}

	resources := make([]Resource, 0)
	for _, role := range resp.Users {
		resp, err := svc.ListAttachedUserPolicies(
			&iam.ListAttachedUserPoliciesInput{
				UserName: role.UserName,
			})
		if err != nil {
			return nil, err
		}

		for _, pol := range resp.AttachedPolicies {
			resources = append(resources, &IAMUserPolicyAttachment{
				svc:        svc,
				policyArn:  *pol.PolicyArn,
				policyName: *pol.PolicyName,
				userName:   *role.UserName,
			})
		}
	}

	return resources, nil
}

func (e *IAMUserPolicyAttachment) Remove() error {
	_, err := e.svc.DetachUserPolicy(
		&iam.DetachUserPolicyInput{
			PolicyArn: &e.policyArn,
			UserName:  &e.userName,
		})
	if err != nil {
		return err
	}

	return nil
}

func (e *IAMUserPolicyAttachment) Properties() types.Properties {
	return types.NewProperties().
		Set("PolicyArn", e.policyArn).
		Set("PolicyName", e.policyName).
		Set("UserName", e.userName)
}

func (e *IAMUserPolicyAttachment) String() string {
	return fmt.Sprintf("%s -> %s", e.userName, e.policyName)
}
