// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract FraudProof {
    struct Validator {
        uint256 stake;
        bool isActive;
    }

    struct StateUpdate {
        bytes32 stateHash;
        address submittedBy;
        uint256 submissionTime;
        bool isCritical;
        bool isChallenged;
    }

    address public admin;
    uint256 public challengePeriod = 24 hours; // Challenge period for state updates
    uint256 public criticalUpdateDelay = 24 hours; // Delay for critical updates

    mapping(address => Validator) public validators;
    mapping(bytes32 => StateUpdate) public stateUpdates;

    event StateSubmitted(bytes32 indexed stateHash, address indexed validator, bool isCritical);
    event StateChallenged(bytes32 indexed stateHash, address indexed challenger);
    event ChallengeResolved(bytes32 indexed stateHash, bool isFraudulent, address slashedValidator);
    event StateApplied(bytes32 indexed stateHash);

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin can perform this action");
        _;
    }

    modifier onlyValidator() {
        require(validators[msg.sender].isActive, "Only active validators can perform this action");
        _;
    }

    constructor() {
        admin = msg.sender;
    }

    // Add a validator
    function addValidator(address _validator, uint256 _stake) external onlyAdmin {
        validators[_validator] = Validator({stake: _stake, isActive: true});
    }

    // Remove a validator
    function removeValidator(address _validator) external onlyAdmin {
        validators[_validator].isActive = false;
    }

    // Submit a state update
    function submitStateUpdate(bytes32 _stateHash, bool _isCritical) external onlyValidator {
        require(stateUpdates[_stateHash].submissionTime == 0, "State update already submitted");

        stateUpdates[_stateHash] = StateUpdate({
            stateHash: _stateHash,
            submittedBy: msg.sender,
            submissionTime: block.timestamp,
            isCritical: _isCritical,
            isChallenged: false
        });

        emit StateSubmitted(_stateHash, msg.sender, _isCritical);
    }

    // Challenge a state update
    function challengeStateUpdate(bytes32 _stateHash) external onlyValidator {
        StateUpdate storage stateUpdate = stateUpdates[_stateHash];
        require(stateUpdate.submissionTime > 0, "State update does not exist");
        require(!stateUpdate.isChallenged, "State update already challenged");
        require(
            block.timestamp <= stateUpdate.submissionTime + challengePeriod,
            "Challenge period has ended"
        );

        stateUpdate.isChallenged = true;

        emit StateChallenged(_stateHash, msg.sender);
    }

    // Resolve a challenge
    function resolveChallenge(bytes32 _stateHash, bool isFraudulent) external onlyAdmin {
        StateUpdate storage stateUpdate = stateUpdates[_stateHash];
        require(stateUpdate.submissionTime > 0, "State update does not exist");
        require(stateUpdate.isChallenged, "State update not challenged");

        if (isFraudulent) {
            address validator = stateUpdate.submittedBy;
            uint256 slashedStake = validators[validator].stake / 2; // Slash 50% of the validator's stake
            validators[validator].stake -= slashedStake;
            if (validators[validator].stake == 0) {
                validators[validator].isActive = false; // Deactivate validator if stake is zero
            }
            delete stateUpdates[_stateHash];
            emit ChallengeResolved(_stateHash, true, validator);
        } else {
            stateUpdate.isChallenged = false; // Mark as valid
            emit ChallengeResolved(_stateHash, false, address(0));
        }
    }

    // Apply a critical state update after delay
    function applyCriticalUpdate(bytes32 _stateHash) external onlyAdmin {
        StateUpdate storage stateUpdate = stateUpdates[_stateHash];
        require(stateUpdate.submissionTime > 0, "State update does not exist");
        require(stateUpdate.isCritical, "State update is not critical");
        require(!stateUpdate.isChallenged, "State update is challenged");
        require(
            block.timestamp >= stateUpdate.submissionTime + criticalUpdateDelay,
            "Critical update delay not passed"
        );

        delete stateUpdates[_stateHash];
        emit StateApplied(_stateHash);
    }
}
